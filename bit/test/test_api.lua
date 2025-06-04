local bit = require("bit")

function Test_and(t)
    local result, err = bit.band(1, 0)
    assert(not err, err)
    assert(result == 0, "and get: " .. tostring(result))
    result, err = bit.band(5, 6)
    assert(not err, err)
    assert(result == 4, "and get: " .. tostring(result))
end

function Test_or(t)
    local result, err = bit.bor(1, 0)
    assert(not err, err)
    assert(result == 1, "or get: " .. tostring(result))
    result, err = bit.bor(5, 6)
    assert(not err, err)
    assert(result == 7, "or get: " .. tostring(result))
end

function Test_xor(t)
    local result, err = bit.bxor(1, 0)
    assert(not err, err)
    assert(result == 1, "xor get: " .. tostring(result))
    result, err = bit.bxor(5, 6)
    assert(not err, err)
    assert(result == 3, "xor get: " .. tostring(result))
end

function Test_left_shift(t)
    local result, err = bit.lshift(1, 0)
    assert(not err, err)
    assert(result == 1, "left_shift get: " .. tostring(result))
    result, err = bit.lshift(0xFF, 8)
    assert(not err, err)
    assert(result == 65280, "left_shift get: " .. tostring(result))
end

function Test_right_shift(t)
    local result, err = bit.rshift(42, 2)
    assert(not err, err)
    assert(result == 10, "right_shift get: " .. tostring(result))
    result, err = bit.rshift(0xFF, 4)
    assert(not err, err)
    assert(result == 15, "right_shift get: " .. tostring(result))
end

function Test_not(t)
    local result, err = bit.bnot(65536)
    assert(not err, err)
    assert(result == 4294901759, "not get: " .. tostring(result))
    result, err = bit.bnot(4294901759)
    assert(not err, err)
    assert(result == 65536, "not get: " .. tostring(result))
end
