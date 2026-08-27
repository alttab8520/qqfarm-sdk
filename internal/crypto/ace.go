package crypto

import (
	"encoding/binary"
	"fmt"
)

// Engine is the long-session ACE surface. Runtime implements it; Identity does not.
type Engine interface {
	Heartbeat() error
	ProcessReceived() error
	SendStatus() error
	DetectSpeedHack(elapsedMs int) error
	PullAnti() ([]byte, error)
	PushAnti(data []byte) error
}

type noopEngine struct{}

func (noopEngine) Heartbeat() error          { return nil }
func (noopEngine) ProcessReceived() error    { return nil }
func (noopEngine) SendStatus() error         { return nil }
func (noopEngine) DetectSpeedHack(int) error { return nil }
func (noopEngine) PullAnti() ([]byte, error) { return nil, nil }
func (noopEngine) PushAnti([]byte) error     { return nil }

func AsEngine(enc Encryptor) Engine {
	if e, ok := enc.(Engine); ok {
		return e
	}
	return noopEngine{}
}

func (rt *Runtime) Heartbeat() error {
	return rt.aceCall(expHeartbeat)
}

func (rt *Runtime) ProcessReceived() error {
	return rt.aceCall(expProcess)
}

func (rt *Runtime) SendStatus() error {
	return rt.aceCall(expSendStatus)
}

func (rt *Runtime) DetectSpeedHack(elapsedMs int) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if rt.closed {
		return fmt.Errorf("加密运行时已关闭")
	}
	if rt.mod == nil || rt.mod.ExportedFunction(expSpeedHack) == nil {
		return nil
	}
	if elapsedMs < 0 {
		elapsedMs = 0
	}
	_, err := rt.call(expSpeedHack, uint64(uint32(elapsedMs)))
	return err
}

func (rt *Runtime) PullAnti() ([]byte, error) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if rt.closed {
		return nil, fmt.Errorf("加密运行时已关闭")
	}
	lenPtr, err := rt.alloc(4)
	if err != nil {
		return nil, err
	}
	defer rt.free(lenPtr)
	mem, err := rt.memory()
	if err != nil {
		return nil, err
	}
	if !mem.Write(lenPtr, []byte{0, 0, 0, 0}) {
		return nil, fmt.Errorf("写入长度指针失败")
	}
	dataPtr, err := rt.call(expPullAnti, uint64(lenPtr))
	if err != nil {
		return nil, err
	}
	raw, ok := mem.Read(lenPtr, 4)
	if !ok {
		return nil, fmt.Errorf("读取 AntiData 长度失败")
	}
	n := int32(binary.LittleEndian.Uint32(raw))
	if dataPtr == 0 || n <= 0 {
		return nil, nil
	}
	out, ok := mem.Read(uint32(dataPtr), uint32(n))
	if !ok {
		return nil, fmt.Errorf("读取 AntiData 失败")
	}
	cp := make([]byte, len(out))
	copy(cp, out)
	return cp, nil
}

func (rt *Runtime) PushAnti(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if rt.closed {
		return fmt.Errorf("加密运行时已关闭")
	}
	ptr, err := rt.alloc(len(data))
	if err != nil {
		return err
	}
	defer rt.free(ptr)
	mem, err := rt.memory()
	if err != nil {
		return err
	}
	if !mem.Write(ptr, data) {
		return fmt.Errorf("写入 AntiData 回灌失败")
	}
	_, err = rt.call(expPushAnti, uint64(ptr), uint64(uint32(len(data))))
	return err
}

func (rt *Runtime) aceCall(name string) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if rt.closed {
		return fmt.Errorf("加密运行时已关闭")
	}
	_, err := rt.call(name)
	return err
}
