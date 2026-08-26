package crypto

import (
	"context"
	"encoding/binary"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

var monoStart = time.Now()

var officialRuntimeTable = []byte{
	93, 86, 110, 34, 65, 129, 8, 113, 53, 192, 121, 32, 86, 162, 255, 139,
	217, 70, 223, 0, 45, 176, 85, 103, 234, 116, 120, 194, 206, 7, 176, 222,
	56, 6, 161, 159, 154, 231, 93, 229, 39, 107, 197, 136, 167, 52, 155, 228,
	209, 117, 218, 8, 107, 241, 32, 62, 53, 200, 238,
}

func recoverQuiet() {
	_ = recover()
}

func (rt *Runtime) instantiateHost(r wazero.Runtime) error {
	_, err := r.NewHostModuleBuilder("a").
		NewFunctionBuilder().WithFunc(rt.impAssert).Export("a").
		NewFunctionBuilder().WithFunc(rt.impWriteFile).Export("b").
		NewFunctionBuilder().WithFunc(rt.impWriteStack).Export("c").
		NewFunctionBuilder().WithFunc(rt.impVersion).Export("d").
		NewFunctionBuilder().WithFunc(rt.impAceVM).Export("e").
		NewFunctionBuilder().WithFunc(rt.impTouch).Export("f").
		NewFunctionBuilder().WithFunc(rt.impReadFile).Export("g").
		NewFunctionBuilder().WithFunc(rt.impClock).Export("h").
		NewFunctionBuilder().WithFunc(rt.impDataDir).Export("i").
		NewFunctionBuilder().WithFunc(rt.impDevice).Export("j").
		NewFunctionBuilder().WithFunc(rt.impFeatureTable).Export("k").
		NewFunctionBuilder().WithFunc(rt.impPlatform).Export("l").
		NewFunctionBuilder().WithFunc(rt.impAppID).Export("m").
		NewFunctionBuilder().WithFunc(rt.impAppID).Export("n").
		NewFunctionBuilder().WithFunc(rt.impFuncJS).Export("o").
		NewFunctionBuilder().WithFunc(rt.impStat).Export("p").
		NewFunctionBuilder().WithFunc(rt.impServerTime).Export("q").
		NewFunctionBuilder().WithFunc(rt.impGrow).Export("r").
		NewFunctionBuilder().WithFunc(rt.impNow).Export("s").
		NewFunctionBuilder().WithFunc(rt.impAppendFile).Export("t").
		NewFunctionBuilder().WithFunc(rt.impAbort).Export("u").
		NewFunctionBuilder().WithFunc(rt.impTQOS).Export("v").
		Instantiate(rt.ctx)
	return err
}

func (rt *Runtime) impAssert(ctx context.Context, mod api.Module, expr, file, line, fn int32) {
	defer recoverQuiet()
}

func (rt *Runtime) impWriteFile(ctx context.Context, mod api.Module, filePtr, dataPtr, encPtr int32) (ret int32) {
	defer recoverQuiet()
	target, err := rt.resolveDataPath(mod, filePtr)
	if err != nil {
		return 0
	}
	text, err := rt.readCString(mod, uint32(dataPtr), 1<<20)
	if err != nil {
		return 0
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return 0
	}
	if err := os.WriteFile(target, []byte(text), 0o644); err != nil {
		return 0
	}
	return 1
}

func (rt *Runtime) impWriteStack(ctx context.Context, mod api.Module, ptr, capacity int32) (ret int32) {
	defer recoverQuiet()
	if ptr == 0 || capacity <= 1 {
		return 0
	}
	stack := debug.Stack()
	room := int(capacity) - 1
	if len(stack) > room {
		stack = stack[len(stack)-room:]
	}
	out := append(append([]byte{}, stack...), 0)
	mem := moduleMemory(mod)
	if mem == nil || !mem.Write(uint32(ptr), out) {
		return 0
	}
	return int32(len(out))
}

func (rt *Runtime) impVersion(ctx context.Context, mod api.Module, ptr, cap int32) (ret int32) {
	defer recoverQuiet()
	return rt.writeCString(mod, rt.version, ptr, cap)
}

func (rt *Runtime) impAceVM(ctx context.Context, mod api.Module, arg int32) (ret int32) {
	defer recoverQuiet()
	return 0
}

func (rt *Runtime) impTouch(ctx context.Context, mod api.Module) {
	defer recoverQuiet()
}

func (rt *Runtime) impReadFile(ctx context.Context, mod api.Module, filePtr, outPtr, capacity, encPtr int32) (ret int32) {
	defer recoverQuiet()
	target, err := rt.resolveDataPath(mod, filePtr)
	if err != nil {
		return 0
	}
	raw, err := os.ReadFile(target)
	if err != nil {
		return 0
	}
	return rt.writeCString(mod, string(raw), outPtr, capacity)
}

func (rt *Runtime) impClock(ctx context.Context, mod api.Module, clockID, lo, hi, outPtr int32) (ret int32) {
	defer recoverQuiet()
	if clockID < 0 || clockID > 3 {
		return 28
	}
	var nanos uint64
	if clockID == 0 {
		nanos = uint64(time.Now().UnixNano())
	} else {
		nanos = uint64(time.Since(monoStart).Nanoseconds())
	}
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, nanos)
	mem := moduleMemory(mod)
	if mem == nil || !mem.Write(uint32(outPtr), buf) {
		return 28
	}
	return 0
}

func (rt *Runtime) impDataDir(ctx context.Context, mod api.Module, ptr, cap int32) (ret int32) {
	defer recoverQuiet()
	dir := filepath.ToSlash(rt.dataDir)
	if dir != "" && !stringsHasSlashSuffix(dir) {
		dir += "/"
	}
	return rt.writeCString(mod, dir, ptr, cap)
}

func (rt *Runtime) impDevice(ctx context.Context, mod api.Module, ptr, cap int32) (ret int32) {
	defer recoverQuiet()
	return rt.writeCString(mod, rt.device, ptr, cap)
}

func (rt *Runtime) impFeatureTable(ctx context.Context, mod api.Module, ptr, cap int32) (ret int32) {
	defer recoverQuiet()
	return rt.writeBytes(mod, officialRuntimeTable, ptr, cap)
}

func (rt *Runtime) impPlatform(ctx context.Context) (ret int32) {
	defer recoverQuiet()
	return 2
}

func (rt *Runtime) impAppID(ctx context.Context, mod api.Module, ptr, cap int32) (ret int32) {
	defer recoverQuiet()
	return rt.writeCString(mod, rt.appID, ptr, cap)
}

func (rt *Runtime) impFuncJS(ctx context.Context, mod api.Module, a1, a2, a3, a4 int32) {
	defer recoverQuiet()
}

func (rt *Runtime) impStat(ctx context.Context, mod api.Module, filePtr int32) (ret int32) {
	defer recoverQuiet()
	return 0
}

func (rt *Runtime) impServerTime(ctx context.Context, mod api.Module, outPtr int32) (ret int32) {
	defer recoverQuiet()
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(int32(time.Now().Unix())))
	mem := moduleMemory(mod)
	if mem == nil || !mem.Write(uint32(outPtr), buf) {
		return 0
	}
	return 1
}

func (rt *Runtime) impGrow(ctx context.Context, mod api.Module, size int32) (ret int32) {
	defer recoverQuiet()
	return 0
}

func (rt *Runtime) impNow(ctx context.Context) (ret float64) {
	defer recoverQuiet()
	return float64(time.Now().UnixMilli())
}

func (rt *Runtime) impAppendFile(ctx context.Context, mod api.Module, filePtr, dataPtr, encPtr int32) (ret int32) {
	defer recoverQuiet()
	target, err := rt.resolveDataPath(mod, filePtr)
	if err != nil {
		return 0
	}
	text, err := rt.readCString(mod, uint32(dataPtr), 1<<20)
	if err != nil {
		return 0
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return 0
	}
	f, err := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return 0
	}
	_, err = f.WriteString(text)
	_ = f.Close()
	if err != nil {
		return 0
	}
	return 1
}

func (rt *Runtime) impAbort(ctx context.Context, mod api.Module) {
	defer recoverQuiet()
}

func (rt *Runtime) impTQOS(ctx context.Context, mod api.Module, ptr, length int32) (ret int32) {
	defer recoverQuiet()
	if length <= 0 {
		return 0
	}
	if _, ok := rt.readBytes(mod, uint32(ptr), uint32(length)); !ok {
		return 0
	}
	return 1
}

func stringsHasSlashSuffix(s string) bool {
	return len(s) > 0 && (s[len(s)-1] == '/' || s[len(s)-1] == '\\')
}
