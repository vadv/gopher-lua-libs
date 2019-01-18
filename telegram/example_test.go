package telegram_test

import (
	"log"

	http "github.com/vadv/gopher-lua-libs/http"
	inspect "github.com/vadv/gopher-lua-libs/inspect"
	telegram "github.com/vadv/gopher-lua-libs/telegram"
	lua "github.com/yuin/gopher-lua"
)

// example sendMessage: https://core.telegram.org/bots/api#sendmessage
func ExampleSendMessage() {
	state := lua.NewState()
	telegram.Preload(state)
	http.Preload(state)
	inspect.Preload(state)
	source := `
local bot = telegram.bot("token")
bot:sendMessage({
    chat_id = number, -- Unique identifier for the target chat
    chat_id or username = "", -- Username of the target channel (in the format @channelusername)
    text = "", -- Text of the message to be sent
    parse_mode = "markdown|html", -- Send Markdown or HTML, if you want Telegram apps to show bold, italic, fixed-width text or inline URLs in your bot's message.
    disable_web_page_preview = false|true, -- Disables link previews for links in this message
    disable_notification = false|true, -- Sends the message silently. Users will receive a notification with no sound.
    reply_to_message_id = nil|number, -- If the message is a reply, ID of the original message
    reply_markup = nil | table, -- https://core.telegram.org/bots/api#inlinekeyboardmarkup, https://core.telegram.org/bots/api#replykeyboardmarkup, https://core.telegram.org/bots/api#replykeyboardremove, https://core.telegram.org/bots/api#forcereply
})
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
}

// example forwardMessage: https://core.telegram.org/bots/api#forwardmessage
func ExampleForwardMessage() {
	state := lua.NewState()
	telegram.Preload(state)
	source := `
local bot = telegram.bot("token")
bot:forwardMessage({
    chat_id = number, -- Unique identifier for the target chat
    chat_id or username = "", -- Username of the target channel (in the format @channelusername)
    from_chat_id = number, -- Unique identifier for the chat where the original message was sent
    disable_notification = false|true, -- Sends the message silently. Users will receive a notification with no sound.
    message_id = number, -- Message identifier in the chat specified in from_chat_id
})
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
}

// example sendPhoto: https://core.telegram.org/bots/api#sendphoto
func ExampleSendPhoto() {
	state := lua.NewState()
	telegram.Preload(state)
	source := `
local bot = telegram.bot("token")
bot:sendPhoto({
    chat_id = number, -- Unique identifier for the target chat
    chat_id or username = "", -- Username of the target channel (in the format @channelusername)
    photo = "path/to/filename|inputfile", -- Photo to send. Pass a file_id as String to send a photo that exists on the Telegram servers (recommended), pass an HTTP URL as a String for Telegram to get a photo from the Internet, or upload a new photo. https://core.telegram.org/bots/api#inputfile
    caption = "", -- Photo caption (may also be used when resending photos by file_id), 0-1024 characters
    parse_mode = "markdown|html", -- Send Markdown or HTML, if you want Telegram apps to show bold, italic, fixed-width text or inline URLs in your bot's message.
    disable_notification = false|true, -- Sends the message silently. Users will receive a notification with no sound.
    reply_to_message_id = nil|number, -- If the message is a reply, ID of the original message
    reply_markup = nil | table, -- https://core.telegram.org/bots/api#inlinekeyboardmarkup, https://core.telegram.org/bots/api#replykeyboardmarkup, https://core.telegram.org/bots/api#replykeyboardremove, https://core.telegram.org/bots/api#forcereply
})
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
}

// example getUpdates: https://core.telegram.org/bots/api#getupdates
func ExampleGetUpdates() {
	state := lua.NewState()
	telegram.Preload(state)
	source := `
local bot = telegram.bot("token")
local updates, err = bot:getUpdates() -- auto offset
for _, update in pairs(updates) do
    inspect(update) -- https://core.telegram.org/bots/api#message
end
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
}

// example deleteMessage: https://core.telegram.org/bots/api#deletemessage
func ExampleDeleteMessage() {
	state := lua.NewState()
	telegram.Preload(state)
	source := `
local bot = telegram.bot("token")
bot:forwardMessage({
    chat_id = number, -- Unique identifier for the target chat
    message_id = number, -- Message identifier in the chat specified in from_chat_id
})
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
}
