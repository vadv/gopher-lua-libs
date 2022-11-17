local base64 = require 'base64'
local strings = require 'strings'
local assert = require 'assert'

function TestEncodeToString(t)
    local tests = {
        {
            name = "input with \1 chars and RawStdEncoding",
            input = "foo\01bar",
            encoder = base64.RawStdEncoding,
            expected = "Zm9vAWJhcg",
        },
        {
            name = "input with \1 chars and StdEncoding",
            input = "foo\01bar",
            encoder = base64.StdEncoding,
            expected = "Zm9vAWJhcg==",
        },
        {
            name = "input with <> chars and RawURLEncoding",
            input = "this is a <tag> and should be encoded",
            encoder = base64.RawURLEncoding,
            expected = "dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA",
        },
        {
            name = "input with <> chars and URLEncoding",
            input = "this is a <tag> and should be encoded",
            encoder = base64.URLEncoding,
            expected = "dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA==",
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            local got = tt.encoder:encode_to_string(tt.input)
            assert:Equal(t, tt.expected, got)
        end)
    end
end

function TestDecodeString(t)
    local tests = {
        {
            name = "input with \1 chars and RawStdEncoding",
            input = "Zm9vAWJhcg",
            encoder = base64.RawStdEncoding,
            expected = "foo\01bar",
        },
        {
            name = "input with \1 chars and StdEncoding",
            input = "Zm9vAWJhcg==",
            encoder = base64.StdEncoding,
            expected = "foo\01bar",
        },
        {
            name = "input with <> chars and RawURLEncoding",
            input = "dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA",
            encoder = base64.RawURLEncoding,
            expected = "this is a <tag> and should be encoded",
        },
        {
            name = "input with <> chars and URLEncoding",
            input = "dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA==",
            encoder = base64.URLEncoding,
            expected = "this is a <tag> and should be encoded",
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            local got, err = tt.encoder:decode_string(tt.input)
            if tt.want_err then
                assert:Error(t, err)
                return
            end
            assert:NoError(t, err)
            assert:Equal(t, tt.expected, got)
        end)
    end
end

function TestEncodeDecode(t)
    local tests = {
        {
            name = "input with \1 chars and RawStdEncoding",
            input = "foo\01bar",
            encoder = base64.RawStdEncoding,
        },
        {
            name = "input with \1 chars and StdEncoding",
            input = "foo\01bar",
            encoder = base64.StdEncoding,
        },
        {
            name = "input with <> chars and RawURLEncoding",
            input = "this is a <tag> and should be encoded",
            encoder = base64.RawURLEncoding,
        },
        {
            name = "input with <> chars and URLEncoding",
            input = "this is a <tag> and should be encoded",
            encoder = base64.URLEncoding,
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            local encoded = tt.encoder:encode_to_string(tt.input)
            local decoded, err = tt.encoder:decode_string(encoded)
            assert:NoError(t, err)
            assert:Equal(t, tt.input, decoded)
        end)
    end
end

function TestEncoder(t)
    local writer = strings.new_builder()
    local encoder = base64.new_encoder(base64.StdEncoding, writer)
    encoder:write("foo", "bar", "baz")
    encoder:close()
    local s = writer:string()
    assert:Equal(t, "Zm9vYmFyYmF6", s)
end

function TestDecoder(t)
    local reader = strings.new_reader("Zm9vYmFyYmF6")
    local decoder = base64.new_decoder(base64.StdEncoding, reader)
    local s = decoder:read("*a")
    assert:Equal(t, "foobarbaz", s)
end

function TestDecoderReadNum(t)
    local encoded = base64.StdEncoding:encode_to_string("123 456 789")
    local reader = strings.new_reader(encoded)
    local decoder = base64.new_decoder(base64.StdEncoding, reader)
    local n = decoder:read("*n")
    assert:Equal(t, 123, n)
    n = decoder:read("*n")
    assert:Equal(t, 456, n)
    n = decoder:read("*n")
    assert:Equal(t, 789, n)
end

function TestDecoderReadCount(t)
    local encoded = base64.StdEncoding:encode_to_string("123 456 789")
    local reader = strings.new_reader(encoded)
    local decoder = base64.new_decoder(base64.StdEncoding, reader)
    local s = decoder:read(3)
    assert:Equal(t, "123", s)
end

function TestDecoderReadline(t)
    local encoded = base64.StdEncoding:encode_to_string("foo\nbar")
    local reader = strings.new_reader(encoded)
    local decoder = base64.new_decoder(base64.StdEncoding, reader)
    local s = decoder:read("*l")
    assert:Equal(t, "foo", s)
    s = decoder:read("*l")
    assert:Equal(t, "bar", s)
end
