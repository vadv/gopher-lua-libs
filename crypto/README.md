# crypto [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/crypto?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/crypto)

## Functions
- `md5(string)` - return md5 checksum from string.
- `sha256(string)` - return sha256 checksum from string.
- `aes_encrypt(string, string, string, string)` - return AES encrypted hex-encoded ciphertext
- `aes_decrypt(string, string, string, string)` - return AES decrypted hex-encoded plain text

AES support 3 modes: GCM, CBC, and CTR - first parameter is mode, second is hex-encoded key, third is hex-encoded initialization vector or nonce - depending on the mode, and forth is hex-encoded plain text or ciphertext.

## Examples

```lua
local crypto = require("crypto")

-- md5
if not(crypto.md5("1\n") == "b026324c6904b2a9cb4b88d6d61c81d1") then
    error("md5")
end

-- sha256
if not(crypto.sha256("1\n") == "4355a46b19d348dc2f57c046f8ef63d4538ebb936000f3c9ee954a27460dd865") then
    error("sha256")
end

--- aes encrypt in GCM mode
s, err = crypto.aes_encrypt(1, "86e15cbc1cbf510d8f2e51d4b63a2144", "b6b86d581a991a652158bd10", "48656c6c6f20776f726c64")
if not(s == "7ec4e38508a26abf7b46e8dc90a7299d5144bcf045e460c3ef6b3e") then
    error("encrypt AES")
end
assert(not err, err)

--- aes decrypt in GCM mode
s, err = crypto.aes_decrypt(1, "86e15cbc1cbf510d8f2e51d4b63a2144", "b6b86d581a991a652158bd10", "7ec4e38508a26abf7b46e8dc90a7299d5144bcf045e460c3ef6b3e")
if not(s == "48656c6c6f20776f726c64") then
    error("decrypt AES)
end
assert(not err, err)

```