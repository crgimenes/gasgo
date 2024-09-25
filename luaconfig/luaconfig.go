package luaconfig

import (
	"errors"
	"log"
	"sync"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type LuaConfig struct {
	mutex    sync.Mutex
	LuaState *lua.LState
}

var (
	LuaConf *LuaConfig
)

func Setup() {
	LuaConf = &LuaConfig{
		LuaState: lua.NewState(),
	}
}

var (
	ErrorFunctionNotFound = errors.New("function not found")
	ErrorNotAllowedType   = errors.New("not allowed return type")
)

func fromGoToLua(v any) lua.LValue {
	switch v.(type) {
	case int:
		return lua.LNumber(v.(int))
	case uint:
		return lua.LNumber(v.(uint))
	case int8:
		return lua.LNumber(v.(int8))
	case uint8:
		return lua.LNumber(v.(uint8))
	case int16:
		return lua.LNumber(v.(int16))
	case uint16:
		return lua.LNumber(v.(uint16))
	case int32:
		return lua.LNumber(v.(int32))
	case uint32:
		return lua.LNumber(v.(uint32))
	case int64:
		return lua.LNumber(v.(int64))
	case uint64:
		return lua.LNumber(v.(uint64))
	case float64:
		return lua.LNumber(v.(float64))
	case string:
		return lua.LString(v.(string))
	case bool:
		return lua.LBool(v.(bool))
	case []byte:
		return lua.LString(v.([]byte))
	case time.Duration:
		return lua.LNumber(v.(time.Duration).Nanoseconds())
	case time.Time:
		return lua.LNumber(v.(time.Time).Unix())
	case map[string]interface{}:
		return mapToLuaTable(v.(map[string]interface{}))
	default:
		log.Println("not allowed type")
		return lua.LNil
	}
}

func fromLuaToGo(v lua.LValue) (any, error) {
	switch v.Type() {
	case lua.LTNil:
		return nil, nil
	case lua.LTNumber:
		return float64(v.(lua.LNumber)), nil
	case lua.LTString:
		return string(v.(lua.LString)), nil
	case lua.LTBool:
		return bool(v.(lua.LBool)), nil
	case lua.LTTable:
		m := make(map[string]interface{})
		t := v.(*lua.LTable)
		t.ForEach(func(k, v lua.LValue) {
			m[string(k.(lua.LString))] = v
		})
		return m, nil
	default:
		return nil, ErrorNotAllowedType
	}
}

func mapToLuaTable(m map[string]interface{}) *lua.LTable {

	table := LuaConf.LuaState.NewTable()

	for k, v := range m {
		switch v.(type) {
		case int:
			table.RawSetString(k, lua.LNumber(v.(int)))
		case float64:
			table.RawSetString(k, lua.LNumber(v.(float64)))
		case string:
			table.RawSetString(k, lua.LString(v.(string)))
		case bool:
			table.RawSetString(k, lua.LBool(v.(bool)))
		case []byte:
			table.RawSetString(k, lua.LString(v.([]byte)))
		case time.Duration:
			table.RawSetString(k, lua.LNumber(v.(time.Duration).Nanoseconds()))
		case time.Time:
			table.RawSetString(k, lua.LNumber(v.(time.Time).Unix()))
		case map[string]interface{}:
			table.RawSetString(k, mapToLuaTable(v.(map[string]interface{})))
		default:
			table.RawSetString(k, lua.LNil)
		}
	}

	return table
}

func ExecFunc(name string, args ...interface{}) (any, error) {
	LuaConf.mutex.Lock()
	defer LuaConf.mutex.Unlock()

	fn := LuaConf.LuaState.GetGlobal(name)
	if fn.Type() != lua.LTFunction {
		return nil, ErrorFunctionNotFound
	}

	largs := make([]lua.LValue, len(args))
	for i, arg := range args {
		largs[i] = fromGoToLua(arg)
	}

	err := LuaConf.LuaState.CallByParam(lua.P{
		Fn:      fn,
		NRet:    1,
		Protect: true,
	}, largs...)

	if err != nil {
		return nil, err
	}

	ret := LuaConf.LuaState.Get(-1)
	return fromLuaToGo(ret)
}

func GetVar(varName string) (any, error) {
	LuaConf.mutex.Lock()
	defer LuaConf.mutex.Unlock()

	ret := LuaConf.LuaState.GetGlobal(varName)
	return fromLuaToGo(ret)
}

func FuncExists(name string) bool {
	LuaConf.mutex.Lock()
	defer LuaConf.mutex.Unlock()

	fn := LuaConf.LuaState.GetGlobal(name)
	return fn.Type() == lua.LTFunction
}
