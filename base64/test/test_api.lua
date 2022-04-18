local base64 = require("base64")

function TestEncodeToString(t)
    tests = {
        {
            name="input with \1 chars and RawStdEncoding",
            input="foo\01bar",
            encoder=base64.RawStdEncoding,
            expected="Zm9vAWJhcg",
        },
        {
            name="input with \1 chars and StdEncoding",
            input="foo\01bar",
            encoder=base64.StdEncoding,
            expected="Zm9vAWJhcg==",
        },
        {
            name="input with <> chars and RawURLEncoding",
            input="this is a <tag> and should be encoded",
            encoder=base64.RawURLEncoding,
            expected="dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA",
        },
        {
            name="input with <> chars and URLEncoding",
            input="this is a <tag> and should be encoded",
            encoder=base64.URLEncoding,
            expected="dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA==",
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            got, err = tt.encoder:encode_to_string(tt.input)
            if tt.want_err then
                assert(err, "expected err")
                return
            end
            assert(not err, err)
            assert(tt.expected == got, string.format("'%s' ~= '%s'", tt.expected, got))
        end)
    end
end

function TestDecodeString(t)
    tests = {
        {
            name="input with \1 chars and RawStdEncoding",
            input="Zm9vAWJhcg",
            encoder=base64.RawStdEncoding,
            expected="foo\01bar",
        },
        {
            name="input with \1 chars and StdEncoding",
            input="Zm9vAWJhcg==",
            encoder=base64.StdEncoding,
            expected="foo\01bar",
        },
        {
            name="input with <> chars and RawURLEncoding",
            input="dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA",
            encoder=base64.RawURLEncoding,
            expected="this is a <tag> and should be encoded",
        },
        {
            name="input with <> chars and URLEncoding",
            input="dGhpcyBpcyBhIDx0YWc-IGFuZCBzaG91bGQgYmUgZW5jb2RlZA==",
            encoder=base64.URLEncoding,
            expected="this is a <tag> and should be encoded",
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            got, err = tt.encoder:decode_string(tt.input)
            if tt.want_err then
                assert(err, "expected err")
                return
            end
            assert(not err, err)
            assert(tt.expected == got, string.format("'%s' ~= '%s'", tt.expected, got))
        end)
    end
end

function TestEncodeDecode(t)

    tests = {
        {
            name="input with \1 chars and RawStdEncoding",
            input="foo\01bar",
            encoder=base64.RawStdEncoding,
        },
        {
            name="input with \1 chars and StdEncoding",
            input="foo\01bar",
            encoder=base64.StdEncoding,
        },
        {
            name="input with <> chars and RawURLEncoding",
            input="this is a <tag> and should be encoded",
            encoder=base64.RawURLEncoding,
        },
        {
            name="input with <> chars and URLEncoding",
            input="this is a <tag> and should be encoded",
            encoder=base64.URLEncoding,
        },
    }
    for _, tt in ipairs(tests) do
        t:Run(tt.name, function(t)
            encoded, err = tt.encoder:encode_to_string(tt.input)
            assert(not err, err)
            decoded = tt.encoder:decode_string(encoded)
            assert(tt.input == decoded, string.format("'%s' ~= '%s'", tt.input, decoded))
        end)
    end
end
