package telegram

import (
	lua "github.com/yuin/gopher-lua"
)

// Preload adds telegram to the given Lua state's package.preload table. After it
// has been preloaded, it can be loaded using require:
//
//  local telegram = require("telegram")
func Preload(L *lua.LState) {
	L.PreloadModule("telegram", Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {

	telegram_bot_ud := L.NewTypeMetatable(`telegram_bot_ud`)
	L.SetGlobal(`telegram_bot_ud`, telegram_bot_ud)
	L.SetField(telegram_bot_ud, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"getUpdates":             GetUpdates,             // https://core.telegram.org/bots/api#getupdates
		"getOffset":              GetOffset,              // max offset value in updates
		"sendMessage":            SendMessage,            // https://core.telegram.org/bots/api#sendmessage
		"deleteMessage":          DeleteMessage,          // https://core.telegram.org/bots/api#deletemessage
		"forwardMessage":         ForwardMessage,         // https://core.telegram.org/bots/api#forwardmessage
		"sendPhoto":              SendPhoto,              // https://core.telegram.org/bots/api#sendphoto
		"editMessageText":        EditMessageText,        // https://core.telegram.org/bots/api#editMessageText
		"editMessageCaption":     EditMessageCaption,     // https://core.telegram.org/bots/api#editMessageCaption
		"editMessageReplyMarkup": EditMessageReplyMarkup, // https://core.telegram.org/bots/api#editMessageReplyMarkup
	}))

	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1
}

var api = map[string]lua.LGFunction{
	"bot": NewBot,
}
