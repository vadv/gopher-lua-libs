local assertions = {
    fail_now = false,
}

function assertions:new(o)
    o = o or {}   -- create object if user does not provide one
    setmetatable(o, self)
    self.__index = self
    return o
end

function assertions:__call(...)
    assert(self.call, 'attempt to call a non-function object')
    return self.call(...)
end

function assertions:cleanseString(s)
    s = string.gsub(s, '\n', '\\n')
    s = string.gsub(s, '\t', '\\t')
    s = string.gsub(s, '\r', '\\r')
    return s
end

function assertions:Fail(t, ...)
    t:LogHelper(2, ...)
    if self.fail_now then
        assert(t.FailNow, "First parameter must be t (the testing.T object)")
        t:FailNow()
    else
        assert(t.Fail, "First parameter must be t (the testing.T object)")
        t:Fail()
    end
    return false
end

function assertions:Failf(t, fmt, ...)
    t:LogHelperf(2, fmt, ...)
    if self.fail_now then
        assert(t.FailNow, "First parameter must be t (the testing.T object)")
        t:FailNow()
    else
        assert(t.Fail, "First parameter must be t (the testing.T object)")
        t:Fail()
    end
    return false
end

function assertions:Equal(t, expected, actual, ...)
    if expected == actual then
        return true
    end
    return self:Fail(t, string.format([[

Error: Not equal:
	   expected: "%s"
	   actual  : "%s"
Messages: ]], self:cleanseString(expected), self:cleanseString(actual)), ...)
end

function assertions:Equalf(t, expected, actual, fmt, ...)
    if expected == actual then
        return true
    end
    return self:Failf(t, string.format([[

Error: Not equal:
	   expected: "%s"
	   actual  : "%s"
Messages: %s
]], self:cleanseString(expected), self:cleanseString(actual), fmt), ...)
end

function assertions:NotEqual(t, expected, actual, ...)
    if expected ~= actual then
        return true
    end
    return self:Fail(t, string.format([[

Error: Should not be %s
Messages: ]], expected), ...)
end

function assertions:NotEqualf(t, expected, actual, fmt, ...)
    if expected ~= actual then
        return true
    end
    return self:Failf(t, string.format([[

Error: Should not be %s
Messages: %s
]], expected, fmt), ...)
end

function assertions:True(t, actual, ...)
    if actual then
        return true
    end
    return self:Fail(t, string.format([[

Error: Should be true
Messages: ]]), ...)
end

function assertions:Truef(t, actual, fmt, ...)
    if actual then
        return true
    end
    return self:Failf(t, string.format([[

Error: Should be true
Messages: %s
]], fmt), ...)

end

function assertions:False(t, actual, ...)
    if not actual then
        return true
    end
    return self:Fail(t, string.format([[

Error Should be false
Messages: ]]), ...)
end

function assertions:Falsef(t, actual, fmt, ...)
    if not actual then
        return true
    end
    return self:Failf(t, string.format([[

Error: Should be false
Messages: %s
]], fmt), ...)
end

return assertions
