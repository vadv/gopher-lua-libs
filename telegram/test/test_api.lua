local http = require("http")
local telegram = require("telegram")
local inspect = require("inspect")

-- local client = http.client({proxy="http://192.168.184.28:3128", insecure_ssl=true})
local bot = telegram.bot("770791440:AAFfrX08qwFtj8YIzcvnhuVzAMv88aqMSxE")

local updates, err = bot:getUpdates()
-- if err then error(err) end

for _, upd in pairs(updates) do
    print(inspect(upd))
    if upd.callback_query then
        bot:sendMessage({
            chat_id = upd.callback_query.message.chat.id,
            reply_to_message_id = upd.callback_query.message.message_id,
            text = "callback query data: "..upd.callback_query.data,
        })
    else
        bot:sendMessage({
            chat_id = upd.message.chat.id,
            reply_to_message_id = upd.message.message_id,
            text = "this is a reply!",
        })
    end
end

local reply_markup_message, err = bot:sendMessage({
    chat_id = 80734283,
    text = "do u like panda?",
    reply_markup={
        inline_keyboard = {
            {
                { text="ok", callback_data="1" }, { text="no", callback_data="2" }
            },
            {
                { text="good", callback_data="3" }, { text="bad", callback_data="4" }
            }
        }
    }
})
-- if err then error(err) end

local _, err = bot:sendPhoto({chat_id = 80734283, caption="panda", photo="./test/panda.jpg"})
-- if err then error(err) end

--[[
local _, err = bot:editMessageReplyMarkup({
    chat_id = reply_markup_message.chat.id,
    message_id = reply_markup_message.message_id,
    reply_markup = {
        inline_keyboard = {
            {
                { text="ok", callback_data="1" }, { text="no", callback_data="2" }
            },
            {
                { text="good (1)", callback_data="3" }, { text="bad", callback_data="4" }
            }
        }
    }
})
if err then error(err) end
-- ]]
