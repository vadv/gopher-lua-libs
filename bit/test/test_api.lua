local bit = require 'bit'
local assert = require 'assert'

function TestAnd(t)
    local tests = {
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
        t:Run(tostring(tt.input1).." and "..tostring(tt.input2) , function(t)
            local got = bit.band(tt.input1, tt.input2)
            assert:Equal(t, tt.expected, got)
        end)
    end
end


function TestOr(t)
    local tests = {
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
        t:Run(tostring(tt.input1).." or "..tostring(tt.input2), function(t)
            local got = bit.bor(tt.input1, tt.input2)
            assert:Equal(t, tt.expected, got)
        end)
    end
end



function TestXor(t)
    local tests = {
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
        t:Run(tostring(tt.input1).." xor "..tostring(tt.input2), function(t)
            local got = bit.bxor(tt.input1, tt.input2)
            assert:Equal(t, tt.expected, got)
        end)
    end
end

function TestLShift(t)
    local tests = {
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
        t:Run(tostring(tt.input1).." << "..tostring(tt.input2), function(t)
            local got = bit.lshift(tt.input1, tt.input2)
            assert:Equal(t, tt.expected, got)
        end)
    end
end

function TestRShift(t)
    local tests = {
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
        t:Run(tostring(tt.input1).." >> "..tostring(tt.input2), function(t)
            local got = bit.rshift(tt.input1, tt.input2)
            assert:Equal(t, tt.expected, got)
        end)
    end
end

function TestNot(t)
    local tests = {
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
        t:Run("not "..tostring(tt.input), function(t)
            local got = bit.bnot(tt.input)
            assert:Equal(t, tt.expected, got)
        end)
    end
end
