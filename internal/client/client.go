package client

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/alttab8520/qqfarm-sdk/internal/crypto"
	"github.com/alttab8520/qqfarm-sdk/internal/game"
	"github.com/alttab8520/qqfarm-sdk/internal/gate"
	"github.com/coder/websocket"
)

const wsBase = "wss://gate-obt.nqf.qq.com/prod/ws"

type Client struct {
	enc crypto.Encryptor

	mu        sync.Mutex
	conn      *websocket.Conn
	pending   map[int64]chan gate.Message
	seq       int64
	serverSeq int64
	user      game.User
	loggedIn  bool
	closed    bool
}

func New() game.Session {
	return NewWith(crypto.Identity{})
}

func NewWith(enc crypto.Encryptor) *Client {
	if enc == nil {
		enc = crypto.Identity{}
	}
	return &Client{enc: enc, pending: map[int64]chan gate.Message{}}
}

func (c *Client) Login(ctx context.Context, in game.LoginIn) (game.User, error) {
	if in.Code == "" {
		return game.User{}, fmt.Errorf("code 不能为空")
	}
	_ = c.Close()
	c.mu.Lock()
	c.pending = map[int64]chan gate.Message{}
	c.seq = 0
	c.serverSeq = 0
	c.closed = false
	c.loggedIn = false
	c.mu.Unlock()

	url := fmt.Sprintf("%s?platform=wx&os=Windows&ver=%s&code=%s&openID=%s",
		wsBase, gameVersion, in.Code, in.OpenID)
	conn, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		HTTPHeader: http.Header{"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64)"}},
	})
	if err != nil {
		return game.User{}, fmt.Errorf("连接网关失败: %w", err)
	}
	conn.SetReadLimit(1 << 20)
	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()
	go c.readLoop()

	body, err := c.call(ctx, "User", "Login", encodeLogin())
	if err != nil {
		_ = c.Close()
		return game.User{}, err
	}
	user, err := decodeUser(body)
	if err != nil {
		_ = c.Close()
		return game.User{}, err
	}
	c.mu.Lock()
	c.user = user
	c.loggedIn = true
	c.mu.Unlock()
	return user, nil
}

func (c *Client) Info() (game.User, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.loggedIn {
		return game.User{}, game.ErrNotLogin
	}
	return c.user, nil
}

func (c *Client) Refresh(ctx context.Context) ([]game.Land, error) {
	user, err := c.Info()
	if err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Plant", "AllLands", encodeAllLands(user.GID))
	if err != nil {
		return nil, err
	}
	return decodeLands(body)
}

func (c *Client) Harvest(ctx context.Context, in game.HarvestIn) ([]game.Item, error) {
	user, err := c.Info()
	if err != nil {
		return nil, err
	}
	host := in.HostGID
	if host == 0 {
		host = user.GID
	}
	body, err := c.call(ctx, "Plant", "Harvest", encodeHarvest(in.LandIDs, host, in.IsAll))
	if err != nil {
		return nil, err
	}
	return decodeItems(body)
}

func (c *Client) Plant(ctx context.Context, in game.PlantIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.SeedID <= 0 || len(in.LandIDs) == 0 {
		return fmt.Errorf("seed_id 和 land_ids 不能为空")
	}
	_, err := c.call(ctx, "Plant", "Plant", encodePlant(in.SeedID, in.LandIDs))
	return err
}

func (c *Client) Friends(ctx context.Context) ([]game.Friend, error) {
	if _, err := c.Info(); err != nil {
		return nil, err
	}
	body, err := c.call(ctx, "Friend", "GetAll", nil)
	if err != nil {
		return nil, err
	}
	return decodeFriends(body)
}

func (c *Client) Help(ctx context.Context, in game.HelpIn) error {
	if _, err := c.Info(); err != nil {
		return err
	}
	if in.GID <= 0 {
		return fmt.Errorf("gid 不能为空")
	}
	if _, err := c.call(ctx, "Visit", "Enter", encodeEnter(in.GID)); err != nil {
		return err
	}
	_, err := c.call(ctx, "Plant", "WaterLand", encodeWater(in.LandIDs, in.GID))
	return err
}

func (c *Client) Close() error {
	c.mu.Lock()
	c.closed = true
	conn := c.conn
	c.conn = nil
	c.loggedIn = false
	for seq, ch := range c.pending {
		close(ch)
		delete(c.pending, seq)
	}
	c.mu.Unlock()
	if conn != nil {
		return conn.Close(websocket.StatusNormalClosure, "")
	}
	return nil
}

func (c *Client) call(ctx context.Context, service, method string, body []byte) ([]byte, error) {
	if body == nil {
		body = []byte{}
	}
	sealed, err := c.enc.Encrypt(body)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	if c.conn == nil {
		c.mu.Unlock()
		return nil, fmt.Errorf("未连接网关")
	}
	c.seq++
	seq := c.seq
	serverSeq := c.serverSeq
	ch := make(chan gate.Message, 1)
	c.pending[seq] = ch
	conn := c.conn
	raw := gate.EncodeRequest(service, method, sealed, c.enc.Token(), seq, serverSeq)
	c.mu.Unlock()

	if err := conn.Write(ctx, websocket.MessageBinary, raw); err != nil {
		c.mu.Lock()
		delete(c.pending, seq)
		c.mu.Unlock()
		return nil, err
	}
	timer := time.NewTimer(15 * time.Second)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-timer.C:
		return nil, fmt.Errorf("RPC %s.%s 超时", service, method)
	case msg, ok := <-ch:
		if !ok {
			return nil, fmt.Errorf("连接已关闭")
		}
		if msg.ErrorCode != 0 {
			return nil, fmt.Errorf("RPC %s: code=%d msg=%s", method, msg.ErrorCode, msg.ErrorMessage)
		}
		return msg.Body, nil
	}
}

func (c *Client) readLoop() {
	ctx := context.Background()
	for {
		c.mu.Lock()
		conn := c.conn
		closed := c.closed
		c.mu.Unlock()
		if conn == nil || closed {
			return
		}
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}
		msg, err := gate.Decode(data)
		if err != nil {
			continue
		}
		c.mu.Lock()
		if msg.ServerSeq > c.serverSeq {
			c.serverSeq = msg.ServerSeq
		}
		if msg.MessageType == gate.TypeResponse {
			if ch, ok := c.pending[msg.ClientSeq]; ok {
				delete(c.pending, msg.ClientSeq)
				c.mu.Unlock()
				ch <- msg
				continue
			}
		}
		c.mu.Unlock()
	}
}
