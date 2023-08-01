package lua

import (
	"io"
	"net/http"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestLuaString(t *testing.T) {
	l := lua.NewState()
	defer l.Close()

	if err := l.DoString(`print("hello")`); err != nil {
		t.Fatal(err)
	}
}

func TestLuaScript(t *testing.T) {
	l := lua.NewState()
	defer l.Close()

	if err := l.DoFile("1.lua"); err != nil {
		t.Fatal(err)
	}
}

func TestLuaScript1(t *testing.T) {
	l := lua.NewState()
	defer l.Close()

	if err := l.DoFile("2.lua"); err != nil {
		t.Fatal(err)
	}
}

func TestLuaScript2(t *testing.T) {
	l := lua.NewState()
	defer l.Close()

	err := l.DoFile("3.lua")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLuaModule(t *testing.T) {
	l := lua.NewState()
	defer l.Close()

	l.PreloadModule("http", httpModuleLoader)
	if err := l.DoFile("4.lua"); err != nil {
		t.Fatal(err)
	}
}

func TestLuaUserData(t *testing.T) {
	l := lua.NewState()
	defer l.Close()

	registerPersonType(l)
	if err := l.DoFile("5.lua"); err != nil {
		t.Fatal(err)
	}
}

func httpModuleLoader(l *lua.LState) int {
	mod := l.SetFuncs(l.NewTable(), map[string]lua.LGFunction{
		"get": get,
	})
	mod.RawSetString("name", lua.LString("http"))
	// l.SetField(mod, "name", lua.LString("http"))
	l.Push(mod)
	return 1
}

func get(l *lua.LState) int {
	uri := l.CheckString(1)
	rsp, err := http.Get(uri)
	if err != nil {
		l.Push(lua.LNil)
		l.Push(lua.LString(err.Error()))
		return 2
	}
	defer rsp.Body.Close()

	tab := l.NewTable()
	tab.RawSetString("status", lua.LNumber(rsp.StatusCode))
	headers := l.NewTable()
	for k := range rsp.Header {
		headers.RawSetString(k, lua.LString(rsp.Header.Get(k)))
	}
	tab.RawSetString("headers", headers)
	b, _ := io.ReadAll(rsp.Body)
	tab.RawSetString("body", lua.LString(b))
	l.Push(tab)

	return 1
}

type Person struct {
	Name string
}

const luaPersonTypeName = "person"

// Registers my person type to given L.
func registerPersonType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaPersonTypeName)
	L.SetGlobal("person", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newPerson))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), personMethods))
}

// Constructor
func newPerson(L *lua.LState) int {
	person := &Person{L.CheckString(1)}
	ud := L.NewUserData()
	ud.Value = person
	L.SetMetatable(ud, L.GetTypeMetatable(luaPersonTypeName))
	L.Push(ud)
	return 1
}

// Checks whether the first lua argument is a *LUserData with *Person and returns this *Person.
func checkPerson(L *lua.LState) *Person {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*Person); ok {
		return v
	}
	L.ArgError(1, "person expected")
	return nil
}

var personMethods = map[string]lua.LGFunction{
	"name": personGetSetName,
}

// Getter and setter for the Person#Name
func personGetSetName(L *lua.LState) int {
	p := checkPerson(L)
	if L.GetTop() == 2 {
		p.Name = L.CheckString(2)
		return 0
	}
	L.Push(lua.LString(p.Name))
	return 1
}
