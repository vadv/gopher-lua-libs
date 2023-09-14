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
    -- Protect against nil
    if s == nil then
        return s
    end
    s = tostring(s)
    s = string.gsub(s, '\n', '\\n')
    s = string.gsub(s, '\t', '\\t')
    s = string.gsub(s, '\r', '\\r')
    return s
end

function assertions:Fail(t, ...)
    assert(t.LogHelper, "First parameter must be t (the testing.T object)")
    t:LogHelper(2, ...)
    if self.fail_now then
        t:FailNow()
    else
        t:Fail()
    end
    return false
end

function assertions:Failf(t, fmt, ...)
    assert(t.LogHelperf, "First parameter must be t (the testing.T object)")
    t:LogHelperf(2, fmt, ...)
    if self.fail_now then
        t:FailNow()
    else
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

function assertions:NoError(t, err, ...)
    if not err then
        return true
    end
    return self:Fail(t, string.format([[

Error:      	Received unexpected error:
                %s
Messages: ]], err), ...)
end

function assertions:NoErrorf(t, err, fmt, ...)
    if not err then
        return true
    end
    return self:Fail(t, string.format([[

Error:      	Received unexpected error:
                %s
Messages: %s
]], err, fmt), ...)
end

function assertions:NotNil(t, object, ...)
    if object ~= nil then
        return true
    end
    return self:Fail(t, string.format([[

Error:      	Expected value not to be nil.
Messages: ]]), ...)
end

function assertions:NotNilf(t, object, fmt, ...)
    if object ~= nil then
        return true
    end
    return self:Failf(t, string.format([[

Error:      	Expected value not to be nil.
Messages: %s
]], fmt), ...)
end

function assertions:Nil(t, object, ...)
    if object == nil then
        return true
    end
    return self:Fail(t, string.format([[

Error:      	Expected nil but got "%s"
Messages: ]], object), ...)
end

function assertions:Nilf(t, object, fmt, ...)
    if object == nil then
        return true
    end
    return self:Failf(t, string.format([[

Error:      	Expected nil but got "%s"
Messages: %s
]], object, fmt), ...)
end

function assertions:Error(t, err, ...)
    if err then
        return true
    end
    return self:Failf(t, string.format([[

Error:      	An error is expected but got nil.
Messages: ]]), ...)
end

function assertions:Errorf(t, err, fmt, ...)
    if err then
        return true
    end
    return self:Failf(t, string.format([[

Error:      	An error is expected but got nil.
Messages: %s
]], fmt), ...)
end

function assertions:Greater(t, x, y, ...)
    if x > y then
        return true
    end
    return self:Fail(t, string.format([[

Error:      	"%s" is not greater than "%s"
Messages: ]], tostring(x), tostring(y)), ...)
end

function assertions:Greaterf(t, x, y, fmt, ...)
    if x > y then
        return true
    end
    return self:Failf(t, string.format([[

Error:      	"%s" is not greater than "%s"
Messages: %s
]], tostring(x), tostring(y), fmt), ...)
end

function assertions:GreaterOrEqual(t, x, y, ...)
    if x >= y then
        return true
    end
    return self:Fail(t, string.format([[

Error:      	"%s" is not greater or equal than "%s"
Messages: ]], tostring(x), tostring(y)), ...)
end

function assertions:GreaterOrEqualf(t, x, y, fmt, ...)
    if x >= y then
        return true
    end
    return self:Failf(t, string.format([[

Error:      	"%s" is not greater or equal than "%s"
Messages: %s
]], tostring(x), tostring(y), fmt), ...)
end

function assertions:Less(t, x, y, ...)
    if x < y then
        return true
    end
    return self:Fail(t, string.format([[

Error:      	"%s" is not less than "%s"
Messages: ]], tostring(x), tostring(y)), ...)
end

function assertions:Lessf(t, x, y, fmt, ...)
    if x < y then
        return true
    end
    return self:Failf(t, string.format([[

Error:      	"%s" is not less than "%s"
Messages: %s
]], tostring(x), tostring(y), fmt), ...)
end

function assertions:LessOrEqual(t, x, y, ...)
    if x <= y then
        return true
    end
    return self:Fail(t, string.format([[

Error:      	"%s" is not less or equal than "%s"
Messages: ]], tostring(x), tostring(y)), ...)
end

function assertions:LessOrEqualf(t, x, y, fmt, ...)
    if x <= y then
        return true
    end
    return self:Failf(t, string.format([[

Error:      	"%s" is not less or equal than "%s"
Messages: %s
]], tostring(x), tostring(y), fmt), ...)
end

return assertions
