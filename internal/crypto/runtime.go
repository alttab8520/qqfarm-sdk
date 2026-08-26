package crypto

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

const (
	officialVersion  = "v3.8.2.1783066265"
	officialSHA256   = "705e326caad538d6cccb40cb1bd54573525a42d12215c9da9c9c513ec4850a5f"
	officialByteSize = 160999
	defaultAppID     = "wx5306c5978fdb76e4"
	defaultGameID    = 3167
	defaultAppKey    = "0"
	defaultDevice    = "PC;windows;Windows 10 x64;microsoft;"
	mergedDataKey    = 1871261153
	rawDecryptStr    = "__mergewasm_shared____wasm_decrypt_strings"
	rawStaticInit    = "x"
	expCreateBuffer  = "A"
	expDestroyBuffer = "B"
	expInitRuntime   = "G"
	expInitInfo      = "H"
	expEncrypt       = "ba"
	expMemory        = "w"
	expTouchSet      = "ha"
)

var mergedDataSegments = [][2]int{
	{1024, 5541}, {6580, 8989}, {15585, 33}, {15643, 1}, {15655, 21},
	{15701, 1}, {15713, 21}, {15759, 1}, {15771, 30}, {15826, 14},
	{15875, 1}, {15887, 21}, {15933, 1}, {15945, 671}, {16632, 400},
	{17040, 103}, {67371008, 404},
}

var requiredExports = []string{
	expMemory, expCreateBuffer, expDestroyBuffer, "C", expInitRuntime, expInitInfo,
	"M", "N", "O", "P", "E", "K", "aa", expEncrypt, "ca", rawDecryptStr, rawStaticInit,
}

var staticInitChecks = [][2]uint32{
	{17872, 0x3f800000},
	{18072, 18076},
}

// Runtime is the official encrypt host. One instance owns ~64 MiB of wasm memory.
type Runtime struct {
	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	wazero    wazero.Runtime
	mod       api.Module
	dataDir   string
	appID     string
	gameID    int32
	appKey    string
	device    string
	version   string
	initToken string
	closed    bool
}

func NewRuntime(wasmPath string) (*Runtime, error) {
	dataDir := os.Getenv("FARM_TSDK_DIR")
	if dataDir == "" {
		if cache, err := os.UserCacheDir(); err == nil {
			dataDir = filepath.Join(cache, "qqfarm-sdk", "tsdk")
		} else {
			dataDir = filepath.Join(".", "data", "tsdk")
		}
	}
	appID := envOr("FARM_TSDK_APP_ID", defaultAppID)
	appKey := envOr("FARM_TSDK_APP_KEY", defaultAppKey)
	gameID := int32(defaultGameID)
	if v := os.Getenv("FARM_TSDK_GAME_ID"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("FARM_TSDK_GAME_ID 无效: %w", err)
		}
		gameID = int32(n)
	}
	ctx, cancel := context.WithCancel(context.Background())
	rt := &Runtime{
		ctx:     ctx,
		cancel:  cancel,
		dataDir: dataDir,
		appID:   appID,
		gameID:  gameID,
		appKey:  appKey,
		device:  defaultDevice,
		version: officialVersion,
	}
	if err := rt.init(wasmPath); err != nil {
		_ = rt.closeLocked()
		return nil, err
	}
	return rt, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (rt *Runtime) init(wasmPath string) error {
	raw, err := os.ReadFile(wasmPath)
	if err != nil {
		return fmt.Errorf("读取加密运行时失败: %w", err)
	}
	if len(raw) != officialByteSize {
		return fmt.Errorf("加密运行时长度不符: 期望 %d 实际 %d", officialByteSize, len(raw))
	}
	sum := sha256.Sum256(raw)
	if hex.EncodeToString(sum[:]) != officialSHA256 {
		return fmt.Errorf("加密运行时校验失败: 期望 %s 实际 %s", officialSHA256, hex.EncodeToString(sum[:]))
	}
	if err := os.MkdirAll(rt.dataDir, 0o755); err != nil {
		return err
	}

	r := wazero.NewRuntime(rt.ctx)
	rt.wazero = r
	if err := rt.instantiateHost(r); err != nil {
		return fmt.Errorf("宿主导入失败: %w", err)
	}
	compiled, err := r.CompileModule(rt.ctx, raw)
	if err != nil {
		return fmt.Errorf("编译加密运行时失败: %w", err)
	}
	mod, err := r.InstantiateModule(rt.ctx, compiled, wazero.NewModuleConfig().WithStartFunctions())
	if err != nil {
		return fmt.Errorf("实例化加密运行时失败: %w", err)
	}
	rt.mod = mod

	var missing []string
	for _, name := range requiredExports {
		if name == expMemory {
			if mod.Memory() == nil {
				missing = append(missing, name)
			}
			continue
		}
		if mod.ExportedFunction(name) == nil {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("加密运行时导出缺失: %s", strings.Join(missing, ", "))
	}

	decrypt := mod.ExportedFunction(rawDecryptStr)
	for _, seg := range mergedDataSegments {
		ptr, length := uint64(uint32(seg[0])), uint64(uint32(seg[1]))
		if _, err := decrypt.Call(rt.ctx, ptr, length, uint64(uint32(mergedDataKey))); err != nil {
			return fmt.Errorf("解密数据段 %d+%d 失败: %w", seg[0], seg[1], err)
		}
	}
	if _, err := mod.ExportedFunction(rawStaticInit).Call(rt.ctx); err != nil {
		return fmt.Errorf("静态初始化失败: %w", err)
	}

	keyPtr, err := rt.allocCString(rt.appKey)
	if err != nil {
		return err
	}
	if _, err := rt.call(expInitRuntime, uint64(uint32(rt.gameID)), uint64(keyPtr)); err != nil {
		rt.free(keyPtr)
		return fmt.Errorf("init_runtime 失败: %w", err)
	}
	rt.free(keyPtr)

	if err := rt.assertStaticInit(); err != nil {
		return err
	}
	if fn := mod.ExportedFunction(expTouchSet); fn != nil {
		_, _ = fn.Call(rt.ctx, 1, 0, math.Float64bits(1))
	}
	return nil
}

func (rt *Runtime) Seal(body []byte) ([]byte, string, error) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if rt.closed {
		return nil, "", fmt.Errorf("加密运行时已关闭")
	}
	if len(body) == 0 {
		return []byte{}, rt.tokenLocked(), nil
	}
	ptr, err := rt.alloc(len(body))
	if err != nil {
		return nil, "", err
	}
	defer rt.free(ptr)
	mem, err := rt.memory()
	if err != nil {
		return nil, "", err
	}
	if !mem.Write(ptr, body) {
		return nil, "", fmt.Errorf("写入明文失败")
	}
	if _, err := rt.call(expEncrypt, uint64(ptr), uint64(uint32(len(body)))); err != nil {
		return nil, "", fmt.Errorf("加密失败: %w", err)
	}
	out, ok := mem.Read(ptr, uint32(len(body)))
	if !ok {
		return nil, "", fmt.Errorf("读出密文失败")
	}
	sealed := make([]byte, len(out))
	copy(sealed, out)
	return sealed, rt.tokenLocked(), nil
}

func (rt *Runtime) BindUser(openID string) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if rt.closed {
		return fmt.Errorf("加密运行时已关闭")
	}
	openID = strings.TrimSpace(openID)
	if openID == "" {
		return fmt.Errorf("open_id 不能为空")
	}
	ptr, err := rt.allocCString(openID)
	if err != nil {
		return err
	}
	defer rt.free(ptr)
	if _, err := rt.call(expInitRuntime, uint64(uint32(rt.gameID)), uint64(ptr)); err != nil {
		return fmt.Errorf("绑定用户失败: %w", err)
	}
	res, err := rt.call(expInitInfo)
	if err != nil {
		return fmt.Errorf("读取初始化凭据失败: %w", err)
	}
	cred, err := rt.readCString(rt.mod, uint32(res), 64*1024)
	if err != nil {
		return err
	}
	if cred == "" {
		return fmt.Errorf("初始化凭据为空")
	}
	rt.initToken = cred
	return nil
}

func (rt *Runtime) Close() error {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	return rt.closeLocked()
}

func (rt *Runtime) closeLocked() error {
	if rt.closed && rt.wazero == nil {
		return nil
	}
	rt.closed = true
	rt.initToken = ""
	rt.mod = nil
	if rt.cancel != nil {
		rt.cancel()
		rt.cancel = nil
	}
	if rt.wazero != nil {
		err := rt.wazero.Close(context.Background())
		rt.wazero = nil
		return err
	}
	return nil
}

func (rt *Runtime) tokenLocked() string {
	if rt.initToken != "" {
		tok := rt.initToken
		rt.initToken = ""
		return tok
	}
	return NewGatewayToken()
}

func (rt *Runtime) call(name string, args ...uint64) (uint64, error) {
	if rt.mod == nil {
		return 0, fmt.Errorf("wasm 未实例化")
	}
	fn := rt.mod.ExportedFunction(name)
	if fn == nil {
		return 0, fmt.Errorf("缺少导出 %s", name)
	}
	results, err := fn.Call(rt.ctx, args...)
	if err != nil {
		return 0, err
	}
	if len(results) == 0 {
		return 0, nil
	}
	return results[0], nil
}

func (rt *Runtime) memory() (api.Memory, error) {
	if rt.mod == nil {
		return nil, fmt.Errorf("wasm 未实例化")
	}
	mem := rt.mod.Memory()
	if mem == nil {
		return nil, fmt.Errorf("没有线性内存")
	}
	return mem, nil
}

func (rt *Runtime) alloc(n int) (uint32, error) {
	if n < 1 {
		n = 1
	}
	res, err := rt.call(expCreateBuffer, uint64(uint32(n)))
	if err != nil {
		return 0, err
	}
	ptr := uint32(res)
	if ptr == 0 {
		return 0, fmt.Errorf("分配 %d 字节失败", n)
	}
	return ptr, nil
}

func (rt *Runtime) allocCString(value string) (uint32, error) {
	data := append([]byte(value), 0)
	ptr, err := rt.alloc(len(data))
	if err != nil {
		return 0, err
	}
	mem, err := rt.memory()
	if err != nil {
		rt.free(ptr)
		return 0, err
	}
	if !mem.Write(ptr, data) {
		rt.free(ptr)
		return 0, fmt.Errorf("写入字符串失败")
	}
	return ptr, nil
}

func (rt *Runtime) free(ptr uint32) {
	if ptr == 0 {
		return
	}
	_, _ = rt.call(expDestroyBuffer, uint64(ptr))
}

func (rt *Runtime) assertStaticInit() error {
	mem, err := rt.memory()
	if err != nil {
		return err
	}
	for _, check := range staticInitChecks {
		raw, ok := mem.Read(check[0], 4)
		if !ok {
			return fmt.Errorf("静态初始化自检读失败: %d", check[0])
		}
		got := binary.LittleEndian.Uint32(raw)
		if got != check[1] {
			return fmt.Errorf("静态初始化未生效: [%d] 期望 0x%x 实际 0x%x", check[0], check[1], got)
		}
	}
	return nil
}

func (rt *Runtime) readCString(mod api.Module, ptr uint32, max int) (string, error) {
	if ptr == 0 {
		return "", nil
	}
	mem := moduleMemory(mod)
	if mem == nil {
		return "", fmt.Errorf("没有线性内存")
	}
	size := mem.Size()
	if ptr >= size {
		return "", fmt.Errorf("字符串指针越界")
	}
	end := uint32(max)
	if remain := size - ptr; end > remain {
		end = remain
	}
	raw, ok := mem.Read(ptr, end)
	if !ok {
		return "", fmt.Errorf("读取字符串失败")
	}
	for i, b := range raw {
		if b == 0 {
			return string(raw[:i]), nil
		}
	}
	return "", fmt.Errorf("字符串未以 NUL 结尾")
}

func (rt *Runtime) writeCString(mod api.Module, value string, ptr, capacity int32) int32 {
	if ptr == 0 || capacity <= 0 {
		return 0
	}
	data := append([]byte(value), 0)
	if len(data) > int(capacity) {
		return 0
	}
	mem := moduleMemory(mod)
	if mem == nil || !mem.Write(uint32(ptr), data) {
		return 0
	}
	return ptr
}

func (rt *Runtime) writeBytes(mod api.Module, value []byte, ptr, capacity int32) int32 {
	if ptr == 0 || int(capacity) < len(value) {
		return 0
	}
	mem := moduleMemory(mod)
	if mem == nil || !mem.Write(uint32(ptr), value) {
		return 0
	}
	return int32(len(value))
}

func (rt *Runtime) readBytes(mod api.Module, ptr, length uint32) ([]byte, bool) {
	mem := moduleMemory(mod)
	if mem == nil {
		return nil, false
	}
	return mem.Read(ptr, length)
}

func (rt *Runtime) resolveDataPath(mod api.Module, filePtr int32) (string, error) {
	rel, err := rt.readCString(mod, uint32(filePtr), 4096)
	if err != nil {
		return "", err
	}
	rel = strings.ReplaceAll(rel, "\\", "/")
	rel = strings.TrimLeft(rel, "/")
	root, err := filepath.Abs(rt.dataDir)
	if err != nil {
		return "", err
	}
	target, err := filepath.Abs(filepath.Join(root, rel))
	if err != nil {
		return "", err
	}
	sep := string(os.PathSeparator)
	if target != root && !strings.HasPrefix(target, root+sep) {
		return "", fmt.Errorf("路径越狱")
	}
	return target, nil
}

func moduleMemory(mod api.Module) api.Memory {
	if mod == nil {
		return nil
	}
	return mod.Memory()
}
