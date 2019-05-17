// Package zabbix implements zabbix api bot for lua.
package zabbix

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	lua_http "github.com/vadv/gopher-lua-libs/http/client/interface"
	lua_json "github.com/vadv/gopher-lua-libs/json"
	lua "github.com/yuin/gopher-lua"
)

// NewBot lua zabbix.bot(config table, http_ud.client) return zabbix_bot_ud
// config = {
//    url = "http://zabbix.url",
//    user = "user",
//    password = "password",
//    debug = true,
//}
func NewBot(L *lua.LState) int {
	config := L.CheckTable(1)
	bot := &luaBot{}
	if L.GetTop() > 1 {
		// http client
		ud := L.CheckUserData(2)
		if v, ok := ud.Value.(lua_http.LuaHTTPClient); ok {
			bot.client = v
		} else {
			L.ArgError(2, "must be http_client_ud")
		}
	} else {
		bot.client = lua_http.NewPureClient()
	}
	var err error
	config.ForEach(func(k lua.LValue, v lua.LValue) {
		switch k.String() {
		case `url`, `base_url`:
			bot.baseURL = v.String()
		case `api`, `api_url`:
			bot.apiURL = v.String()
		case `login`, `login_url`:
			bot.loginURL = v.String()
		case `chat`, `chat_url`:
			bot.chatURL = v.String()
		case `user`, `username`:
			bot.user = v.String()
		case `passwd`, `password`:
			bot.password = v.String()
		case `debug`:
			bot.debug = (v.String() == `true`)
		default:
			err = fmt.Errorf("unknown config parameter: `%s`", k.String())
		}
	})
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	if err := bot.updateURLs(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	ud := L.NewUserData()
	ud.Value = bot
	L.SetMetatable(ud, L.GetTypeMetatable("zabbix_bot_ud"))
	L.Push(ud)
	if bot.debug {
		log.Printf("[DEBUG] create zabbix bot: %#v\n", bot)
	}
	return 1
}

// Login lua zabbix_bot_ud:login() return error
func Login(L *lua.LState) int {
	b := checkBot(L, 1)
	err := b.login()
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

// Logout lua zabbix_bot_ud:logout() return error
func Logout(L *lua.LState) int {
	b := checkBot(L, 1)
	err := b.logout()
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}

// ApiVersion lua zabbix_bot_ud:api_version() returns (string, err)
func ApiVersion(L *lua.LState) int {
	b := checkBot(L, 1)
	version, err := b.apiVersion()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(version))
	return 1
}

// Request lua zabbix_bot_ud:request(method, params={}) returns (table, err)
func Request(L *lua.LState) int {
	b := checkBot(L, 1)
	method := L.CheckString(2)
	config := L.CheckTable(3)
	var params map[string]interface{}

	// lua table -> byte -> params

	// input to byte
	data, err := lua_json.ValueEncode(config)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	if err := json.Unmarshal(data, &params); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	// send
	response, err := b.sendRequest(method, params)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	// decode response
	byt, isString, err := response.resultToBytes()
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	if b.debug {
		log.Printf("[DEBUG] response body json: %s\n", byt)
	}
	if isString {
		L.Push(lua.LString(string(byt)))
		return 1
	}
	result, err := lua_json.ValueDecode(L, byt)
	if b.debug {
		log.Printf("[DEBUG] response body lua: %s\n", result)
	}
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(result)
	return 1
}

// SaveGraph lua zabbix_bot_ud:save_graph(itemID, filename, {period, width, height}) return err
// default graph settings {period = 3600, width = 500, height = 300}
func SaveGraph(L *lua.LState) int {
	b := checkBot(L, 1)
	itemId, err := strconv.ParseInt(L.CheckAny(2).String(), 10, 64)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	filename := L.CheckString(3)
	period, width, height := int64(60*60), int64(500), int64(300)
	if L.GetTop() > 3 {
		params := L.CheckTable(4)
		params.ForEach(func(k lua.LValue, v lua.LValue) {
			switch k.String() {
			case `period`:
				period, err = strconv.ParseInt(v.String(), 10, 64)
				if err != nil {
					return
				}
			case `width`:
				width, err = strconv.ParseInt(v.String(), 10, 64)
				if err != nil {
					return
				}
			case `height`:
				height, err = strconv.ParseInt(v.String(), 10, 64)
				if err != nil {
					return
				}
			default:
				err = fmt.Errorf("unknown graph settings parameter: `%s`", k.String())
				return
			}
		})
	}
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	err = b.saveGraph(filename, itemId, period, width, height)
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 1
	}
	return 0
}
