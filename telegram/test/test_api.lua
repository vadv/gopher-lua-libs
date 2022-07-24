local http = require("http")
local telegram = require("telegram")
local inspect = require("inspect")

function Test_telegram(t)
    local client_settings = {}
    if not os.getenv("TRAVIS") then
        client_settings = { proxy = "http://192.168.184.28:3128", insecure_ssl = true }
    end
    local bot = telegram.bot("770791440:AAFfrX08qwFtj8YIzcvnhuVzAMv88aqMSxE", http.client(client_settings))

    local updates, err = bot:getUpdates()
    --assert(not err, err)

    if updates then
        -- concurency build
        for _, upd in pairs(updates) do
            if upd.callback_query then
                local msg, err = bot:sendMessage({
                    chat_id = upd.callback_query.message.chat.id,
                    reply_to_message_id = upd.callback_query.message.message_id,
                    text = "callback query data: " .. upd.callback_query.data,
                })
                if not err then
                    bot:deleteMessage({ chat_id = upd.callback_query.message.chat.id, message_id = msg.message_id })
                end
            else
                local msg, err = bot:sendMessage({
                    chat_id = upd.message.chat.id,
                    reply_to_message_id = upd.message.message_id,
                    text = "this is a reply!",
                })
                if not err then
                    bot:deleteMessage({ chat_id = upd.message.chat.id, message_id = msg.message_id })
                end
            end
        end
    end

    local reply_markup_message, err = bot:sendMessage({
        chat_id = 80734283,
        text = "do u like panda?",
        reply_markup = {
            inline_keyboard = {
                {
                    { text = "ok", callback_data = "1" }, { text = "no", callback_data = "2" }
                },
                {
                    { text = "good", callback_data = "3" }, { text = "bad", callback_data = "4" }
                }
            }
        }
    })
    if not err then
        bot:deleteMessage({
            chat_id = 80734283, message_id = reply_markup_message.message_id
        })
    end
    --assert(not err, err)

    local msg, err = bot:sendPhoto({ chat_id = 80734283, caption = "panda", photo = "./test/panda.jpg" })
    if not (err) then
        bot:deleteMessage({
            chat_id = 80734283, message_id = msg.message_id
        })
    end

    --assert(not err, err)

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
end
