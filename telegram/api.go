// Package telegram implements telegram api bot for lua.
package telegram

import (
	"encoding/json"

	lua_http "github.com/vadv/gopher-lua-libs/http/client/interface"
	lua_json "github.com/vadv/gopher-lua-libs/json"
	lua "github.com/yuin/gopher-lua"
)

// NewBot lua telegram.bot(tocken, http_ud.client) return telegram_bot_ud
func NewBot(L *lua.LState) int {
	token := L.CheckString(1)
	bot := &luaBot{token: token}
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
	ud := L.NewUserData()
	ud.Value = bot
	L.SetMetatable(ud, L.GetTypeMetatable("telegram_bot_ud"))
	L.Push(ud)
	return 1
}

// GetUpdates lua telegram_bot_ud:get_updates() returns (table, err)
func GetUpdates(L *lua.LState) int {
	bot := checkBot(L, 1)
	var c *UpdateConfig
	if L.GetTop() > 1 {
		config := L.CheckTable(2)
		// input to byte
		data, err := lua_json.ValueEncode(config)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		// byte to config
		if err := json.Unmarshal(data, c); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
	}
	updates, err := bot.getUpdates(c)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result, err := lua_json.ValueDecode(L, updatesToBytes(updates))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(result)
	return 1
}

// GetOffset lua telegram_bot_ud:getOffset() returns int
func GetOffset(L *lua.LState) int {
	bot := checkBot(L, 1)
	L.Push(lua.LNumber(bot.offset))
	return 1
}

// lua table -> byte -> Chattable
func sendGeneric(L *lua.LState, iface interface{}) int {
	bot := checkBot(L, 1)
	message := L.CheckTable(2)
	c, ok := iface.(Chattable)
	if !ok {
		panic("must be chattable struct")
	}
	// input to byte
	data, err := lua_json.ValueEncode(message)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	// byte to chattable
	if err := json.Unmarshal(data, c); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	// send
	response, err := bot.send(c)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	// decode response
	result, err := lua_json.ValueDecode(L, response.toBytes())
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(result)
	return 1
}

// SendMessage lua telegram_bot_ud:message(table) returns (table, err)
func SendMessage(L *lua.LState) int {
	return sendGeneric(L, &MessageConfig{})
}

// ForwardMessage lua telegram_bot_ud:forward(table) returns (table, err)
func ForwardMessage(L *lua.LState) int {
	return sendGeneric(L, &ForwardConfig{})
}

// SendPhoto lua telegram_bot_ud:photo(table) returns (table, err)
func SendPhoto(L *lua.LState) int {
	return sendGeneric(L, &PhotoConfig{})
}

// EditMessageText lua telegram_bot_ud:editMessageText(table) returns (bool, err)
func EditMessageText(L *lua.LState) int {
	return sendGeneric(L, &EditMessageTextConfig{})
}

// EditMessageCaption lua telegram_bot_ud:editMessageCaption(table) returns (bool, err)
func EditMessageCaption(L *lua.LState) int {
	return sendGeneric(L, &EditMessageCaptionConfig{})
}

// EditMessageReplyMarkup lua telegram_bot_ud:editMessageReplyMarkup(table) returns (bool, err)
func EditMessageReplyMarkup(L *lua.LState) int {
	return sendGeneric(L, &EditMessageReplyMarkupConfig{})
}

// EditMessageReplyMarkup lua telegram_bot_ud:deleteMessage(table) returns (bool, err)
func DeleteMessage(L *lua.LState) int {
	return sendGeneric(L, &DeleteMessageConfig{})
}
