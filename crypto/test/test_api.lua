local crypto = require("crypto")
local assert = require 'assert'
local hex = require 'hex'
local require = require 'require'
local filepath = require 'filepath'
local ioutil = require 'ioutil'

function TestMD5(t)
    local tests = {
        {
            input = "1\n",
            expected = "b026324c6904b2a9cb4b88d6d61c81d1",
        },
        {
            input = "test",
            expected = "098f6bcd4621d373cade4e832627b4f6"
        }
    }
    for _, tt in ipairs(tests) do
        t:Run("md5(" .. tostring(tt.input) .. ")", function(t)
            local got = crypto.md5(tt.input)
            assert:Equal(t, tt.expected, got)
        end)
    end
end

function TestSha256(t)
    local tests = {
        {
            input = "1\n",
            expected = "4355a46b19d348dc2f57c046f8ef63d4538ebb936000f3c9ee954a27460dd865",
        },
        {
            input = "test",
            expected = "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
        }
    }
    for _, tt in ipairs(tests) do
        t:Run("sha256(" .. tostring(tt.input) .. ")", function(t)
            local got = crypto.sha256(tt.input)
            assert:Equal(t, tt.expected, got)
        end)
    end
end

function TestAESEncryptHex(t)
    local tests = {
        {
            data = "48656c6c6f207w76f726c64", -- "Hello world" in hex
            mode = "GCM",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "b6b86d581a991a652158bd10",
            expected = nil,
            err = "failed to decode params: failed to decode data: encoding/hex: invalid byte: U+0077 'w'",
        },
        {
            data = "48656c6c6f20776f726c64", -- "Hello world" in hex
            mode = "GCM",
            key = "86e15cbc1cbf51d8f2e51d4b63a2144",
            init = "b6b86d581a991a652158bd10",
            expected = nil,
            err = "failed to decode params: failed to decode key: encoding/hex: odd length hex string",
        },
        {
            data = "48656c6c6f20776f726c64", -- "Hello world" in hex
            mode = "GCM",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "b6b86d581a991a652158bd10",
            expected = "7ec4e38508a26abf7b46e8dc90a7299d5144bcf045e460c3ef6b3e",
            err = nil,
        },
        {
            data = "48656c6c6f20776f726c64", -- "Hello world" in hex
            mode = "GCM",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "b6b86d581a991a652158bd010211",
            expected = nil,
            err = "failed to encrypt: incorrect GCM nonce size: 14, expected: 12",
        },
        {
            data = "48656c6c6f20776f726c64", -- "Hello world" in hex
            mode = "GCM",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "b6b86d581a991a652158bd010211",
            expected = nil,
            err = "failed to encrypt: incorrect GCM nonce size: 14, expected: 12",
        },
        {
            data = "48656c6c6f20776f726c64", -- "Hello world" in hex
            mode = "cbc",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "068bb92e032884ba8b260fa7d3a80005",
            expected = "dfba6f71cce4d4b76be301b577d9f095",
            err = nil,
        },
        {
            data = "48656c6c6f20776f726c64", -- "Hello world" in hex
            mode = "CBC",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "068bb92e03288884ba8b260fa7d3a80005",
            expected = nil,
            err = "failed to encrypt: invalid IV size: 17, expected: 16",
        },
        {
            data = "48656c6c6f20776f726c64", -- "Hello world" in hex
            mode = "CTR",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "e3057fc2bf103a09a1b2c3d4e5f60718",
            expected = "138434a80bd7dcd9ee8adc",
            err = nil,
        },
        {
            data = "48656c6c6f20776f726c64", -- "Hello world" in hex
            mode = "CTR",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "e3057fc2b9f103a909a1b2c3d4e5f60718",
            expected = nil,
            err = "failed to encrypt: invalid IV size: 17, expected: 16",
        },
    }
    for _, tt in ipairs(tests) do
        t:Run("aes_encrypt in " .. tostring(tt.mode) .. " mode", function(t)
            local got, err = crypto.aes_encrypt_hex(tt.mode, tt.key, tt.init, tt.data)
            assert:Equal(t, tt.expected, got)
            assert:Equal(t, tt.err, err)
        end)
    end
end

function TestAESDecryptHex(t)
    local tests = {
        {
            data = "7ec4e38508a26abf7b46e8dc90a7299d5144bcf045e460c3efwb3e",
            mode = "GCM",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "b6b86d581a991a652158bd10",
            expected = nil,
            err = "failed to decode params: failed to decode data: encoding/hex: invalid byte: U+0077 'w'",
        },
        {
            data = "7ec4e38508a26abf7b46e8dc90a7299d5144bcf045e460c3ef6b3e",
            mode = "GCM",
            key = "86e15cbc1cbf51d8f2e51d4b63a2144",
            init = "b6b86d581a991a652158bd10",
            expected = nil,
            err = "failed to decode params: failed to decode key: encoding/hex: odd length hex string",
        },
        {
            data = "7ec4e38508a26abf7b46e8dc90a7299d5144bcf045e460c3ef6b3e",
            mode = "GCM",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "b6b86d581a991a652158bd10",
            expected = "48656c6c6f20776f726c64", -- "Hello world" in hex
            err = nil,
        },
        {
            data = "7ec4e38508a26abf7b46e8dc90a7299d5144bcf045e460c3ef6b3e",
            mode = "GCM",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "b6b86d581a991a652158bd010211",
            expected = nil,
            err = "failed to decrypt: incorrect GCM nonce size: 14, expected: 12",
        },
        {
            data = "dfba6f71cce4d4b76be301b577d9f095",
            mode = "cbc",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "068bb92e032884ba8b260fa7d3a80005",
            expected = "48656c6c6f20776f726c640505050505", -- "Hello world" + padding in hex
            err = nil,
        },
        {
            data = "138434a80bd7dcd9ee8adc",
            mode = "CTR",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "e3057fc2bf103a09a1b2c3d4e5f60718",
            expected = "48656c6c6f20776f726c64", -- "Hello world" in hex
            err = nil,
        },
    }
    for _, tt in ipairs(tests) do
        t:Run("aes_decrypt in " .. tostring(tt.mode) .. " mode", function(t)
            local got, err = crypto.aes_decrypt_hex(tt.mode, tt.key, tt.init, tt.data)
            assert:Equal(t, tt.expected, got)
            assert:Equal(t, tt.err, err)
        end)
    end
end

function TestAESEncrypt(t)
    tests = {
        {
            data = "48656c6c6f20776f726c64",
            mode = "cbc",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "068bb92e032884ba8b260fa7d3a80005",
            expected = "dfba6f71cce4d4b76be301b577d9f095",
            wantErr = false,
        },
    }
    for _, tt in ipairs(tests) do
        local key, err = hex.decode_string(tt.key)
        require:NoError(t, err)
        local init, err = hex.decode_string(tt.init)
        require:NoError(t, err)
        local data, err = hex.decode_string(tt.data)
        require:NoError(t, err)
        local got, err = crypto.aes_encrypt(tt.mode, key, init, data)
        if tt.wantErr then
            require:Error(t, err)
            return
        end
        require:NoError(t, err)
        got = hex.encode_to_string(got)
        assert:Equal(t, tt.err, err)
    end
end

function TestAESDecrypt(t)
    tests = {
        {
            data = "138434a80bd7dcd9ee8adc",
            mode = "CTR",
            key = "86e15cbc1cbf510d8f2e51d4b63a2144",
            init = "e3057fc2bf103a09a1b2c3d4e5f60718",
            expected = "48656c6c6f20776f726c64", -- "Hello world" in hex
            wantErr = false,
        },
    }
    for _, tt in ipairs(tests) do
        t:Run("aes_decrypt in " .. tostring(tt.mode) .. " mode", function(t)
            local key, err = hex.decode_string(tt.key)
            require:NoError(t, err)
            local init, err = hex.decode_string(tt.init)
            require:NoError(t, err)
            local data, err = hex.decode_string(tt.data)
            require:NoError(t, err)
            local got, err = crypto.aes_decrypt(tt.mode, key, init, data)
            if tt.wantErr then
                require:Error(t, err)
                return
            end
            require:NoError(t, err)
            got, err = hex.encode_to_string(got)
            require:NoError(t, err)
            assert:Equal(t, tt.expected, got)
        end)
    end
end

function TestAESCodecFile(t)
    for i = 1, 1 do
        local data, err = ioutil.read_file(filepath.join("test/data", tostring(i) .. ".data.bin"))
        require:NoError(t, err)
        local expected, err = ioutil.read_file(filepath.join("test/data", tostring(i) .. ".expected.bin"))
        require:NoError(t, err)
        local init, err = ioutil.read_file(filepath.join("test/data", tostring(i) .. ".init.bin"))
        require:NoError(t, err)
        local key, err = ioutil.read_file(filepath.join("test/data", tostring(i) .. ".key.bin"))
        require:NoError(t, err)
        t:Run("TestAESEncryptFile " .. tostring(i), function(t)
            local got, err = crypto.aes_encrypt("CTR", key, init, data)
            require:NoError(t, err)
            assert:Equal(t, expected, got)

            local decrypted, err = crypto.aes_decrypt("CTR", key, init, got)
            t:Logf('data: "%s", decrypted: "%s"', data, decrypted)
            require:NoError(t, err)
            assert:Equal(t, data, decrypted)
        end)
    end
end