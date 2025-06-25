local hex = require 'hex'
local strings = require 'strings'
local assert = require 'assert'

function TestHexCodec(t)
    local tests = {
        {
            encoded = "48656c6c6f20776f726c64", -- "Hello world" in hex
            decoded = "Hello world",
        },
        {
            encoded = "",
            decoded = "",
        },
    }
    for _, tt in ipairs(tests) do
        t:Run("hex.decode_string(" .. tostring(tt.encoded) .. ")", function(t)
            local got = hex.decode_string(tt.encoded)
            assert:Equal(t, tt.decoded, got)
        end)

        t:Run("hex.encode_to_string(" .. tostring(tt.decoded) .. ")", function(t)
            local got = hex.encode_to_string(tt.decoded)
            assert:Equal(t, tt.encoded, got)
        end)
    end
end

function TestEncoder(t)
    local writer = strings.new_builder()
    local encoder = hex.new_encoder(writer)
    encoder:write("foo", "bar", "baz")
    encoder:close()
    local s = writer:string()
    assert:Equal(t, "666f6f62617262617a", s)
end

function TestDecoder(t)
    local reader = strings.new_reader("666f6f62617262617a")
    local decoder = hex.new_decoder(reader)
    local s = decoder:read("*a")
    assert:Equal(t, "foobarbaz", s)
end
