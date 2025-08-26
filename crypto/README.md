# crypto [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/crypto?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/crypto)

## Functions
- `md5(string)` - return md5 checksum from string.
- `sha256(string)` - return sha256 checksum from string.
- `aes_encrypt(string, string, string, string)` - return AES encrypted binary ciphertext
- `aes_decrypt(string, string, string, string)` - return AES decrypted binary text
- `aes_encrypt_hex(string, string, string, string)` - return AES encrypted hex-encoded ciphertext
- `aes_decrypt_hex(string, string, string, string)` - return AES decrypted hex-encoded plain text

AES support 3 modes: GCM, CBC, and CTR - first parameter is mode, second is hex-encoded key, third is hex-encoded
initialization vector or nonce - depending on the mode, and forth is hex-encoded plain text or ciphertext.

Since lua strings are binary safe, you can use any binary data as input and output and, for your convenience, the
library also provides hex-encoded versions of the encrypt and decrypt functions. The first argument (the mode string)
can be one of the following: "GCM", "CBC", or "CTR" (case-insensitive) and is not hex-encoded for the hex variants.

## Examples

```lua
local crypto = require 'crypto'

-- md5
if not(crypto.md5("1\n") == "b026324c6904b2a9cb4b88d6d61c81d1") then
    error("md5")
end

-- sha256
if not(crypto.sha256("1\n") == "4355a46b19d348dc2f57c046f8ef63d4538ebb936000f3c9ee954a27460dd865") then
    error("sha256")
end

--- aes encrypt in GCM mode with hex-encoded data
s, err = crypto.aes_encrypt_hex(crypto.GCM, "86e15cbc1cbf510d8f2e51d4b63a2144", "b6b86d581a991a652158bd10", "48656c6c6f20776f726c64")
assert(not err, err)
if not(s == "7ec4e38508a26abf7b46e8dc90a7299d5144bcf045e460c3ef6b3e") then
    error("encrypt AES")
end

--- aes decrypt in GCM mode with hex-encoded data
s, err = crypto.aes_decrypt_hex(crypto.GCM, "86e15cbc1cbf510d8f2e51d4b63a2144", "b6b86d581a991a652158bd10", "7ec4e38508a26abf7b46e8dc90a7299d5144bcf045e460c3ef6b3e")
assert(not err, err)
if not(s == "48656c6c6f20776f726c64") then
    error("decrypt AES")
end

--- Binary examples setup of binary strings equivalent to the hex-encoded strings above:
local hex = require 'hex'
local key, iv, plaintext, encrypted, err
key, err = hex.decode_string('86e15cbc1cbf510d8f2e51d4b63a2144')
assert(not err, err)
iv, err = hex.decode_string('b6b86d581a991a652158bd10')
assert(not err, err)
plaintext, err = hex.decode_string('48656c6c6f20776f726c64')
assert(not err, err)
s, err = crypto.aes_encrypt(crypto.GCM, key, iv, plaintext)
assert(not err, err)
encrypted, err = hex.decode_string("7ec4e38508a26abf7b46e8dc90a7299d5144bcf045e460c3ef6b3e")
assert(not err, err)

--- aes encrypt binary in GCM mode
s, err = crypto.aes_encrypt(crypto.GCM, key, iv, plaintext)
assert(not err, err)
if not(s == encrypted) then
    error("encrypt AES")
end

--- aes decrypt in GCM mode
s, err = crypto.aes_decrypt(crypto.GCM, key, iv, encrypted)
assert(not err, err)
if not(s == plaintext) then
    error("decrypt AES")
end

```