local bit = require 'bit'
local assert = require 'assert'

function TestAnd(t)
    local tests = {
        {
            input1 = -3,
            input2 = 23,
            expected = nil,
            err = "cannot convert negative int -3 to uint32",
        },
        {
            input1 = 4294967296,
            input2 = 23,
            expected = nil,
            err = "int 4294967296 overflows uint32",
        },
        {
            input1 = 1,
            input2 = 0,
            expected = 0,
        },
        {
            input1 = 111,
            input2 = 222,
            expected = 78,
        }
    }
    for _, tt in ipairs(tests) do
        t:Run(tostring(tt.input1) .. " and " .. tostring(tt.input2), function(t)
            local got, err = bit.band(tt.input1, tt.input2)
            assert:Equal(t, tt.expected, got)
            assert:Equal(t, tt.err, err)
        end)
    end
end

function TestOr(t)
    local tests = {
        {
            input1 = 5,
            input2 = -423,
            expected = nil,
            err = "cannot convert negative int -423 to uint32",
        },
        {
            input1 = 123,
            input2 = 4294967296,
            expected = nil,
            err = "int 4294967296 overflows uint32",
        },
        {
            input1 = 1,
            input2 = 0,
            expected = 1,
        },
        {
            input1 = 111,
            input2 = 222,
            expected = 255,
        }
    }
    for _, tt in ipairs(tests) do
        t:Run(tostring(tt.input1) .. " or " .. tostring(tt.input2), function(t)
            local got, err = bit.bor(tt.input1, tt.input2)
            assert:Equal(t, tt.expected, got)
            assert:Equal(t, tt.err, err)
        end)
    end
end

function TestXor(t)
    local tests = {
        {
            input1 = -1,
            input2 = -46,
            expected = nil,
            err = "cannot convert negative int -1 to uint32",
        },
        {
            input1 = 4294967300,
            input2 = 46,
            expected = nil,
            err = "int 4294967300 overflows uint32",
        },
        {
            input1 = 1,
            input2 = 0,
            expected = 1,
        },
        {
            input1 = 111,
            input2 = 222,
            expected = 177,
        }
    }
    for _, tt in ipairs(tests) do
        t:Run(tostring(tt.input1) .. " xor " .. tostring(tt.input2), function(t)
            local got, err = bit.bxor(tt.input1, tt.input2)
            assert:Equal(t, tt.expected, got)
            assert:Equal(t, tt.err, err)
        end)
    end
end

function TestLShift(t)
    local tests = {
        {
            input1 = 0,
            input2 = -10,
            expected = nil,
            err = "cannot convert negative int -10 to uint32",
        },
        {
            input1 = 4294967297,
            input2 = 4294967298,
            expected = nil,
            err = "int 4294967297 overflows uint32",
        },
        {
            input1 = 123456,
            input2 = 8,
            expected = 31604736,
        },
        {
            input1 = 0XFF,
            input2 = 8,
            expected = 65280,
        }
    }
    for _, tt in ipairs(tests) do
        t:Run(tostring(tt.input1) .. " << " .. tostring(tt.input2), function(t)
            local got, err = bit.lshift(tt.input1, tt.input2)
            assert:Equal(t, tt.expected, got)
            assert:Equal(t, tt.err, err)
        end)
    end
end

function TestRShift(t)
    local tests = {
        {
            input1 = -10,
            input2 = 0,
            expected = nil,
            err = "cannot convert negative int -10 to uint32",
        },
        {
            input1 = 4294967296,
            input2 = -3,
            expected = nil,
            err = "int 4294967296 overflows uint32",
        },
        {
            input1 = 123456,
            input2 = 8,
            expected = 482,
        },
        {
            input1 = 0XFF,
            input2 = 1,
            expected = 0x7F,
        }
    }
    for _, tt in ipairs(tests) do
        t:Run(tostring(tt.input1) .. " >> " .. tostring(tt.input2), function(t)
            local got, err = bit.rshift(tt.input1, tt.input2)
            assert:Equal(t, tt.expected, got)
            assert:Equal(t, tt.err, err)
        end)
    end
end

function TestNot(t)
    local tests = {
        {
            input = -3,
            expected = nil,
            err = "cannot convert negative int -3 to uint32",
        },
        {
            input = 4294967297,
            expected = nil,
            err = "int 4294967297 overflows uint32",
        },
        {
            input = 65536,
            expected = 4294901759,
        },
        {
            input = 4294901759,
            expected = 65536,
        }
    }
    for _, tt in ipairs(tests) do
        t:Run("not " .. tostring(tt.input), function(t)
            local got, err = bit.bnot(tt.input)
            assert:Equal(t, tt.expected, got)
            assert:Equal(t, tt.err, err)
        end)
    end
end
