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

function assertions:Fail(t, ...)
    t:LogHelper(2, ...)
    if self.fail_now then
        print("Failing now")
        t:FailNow()
    else
        print("Failing")
        t:Fail()
    end
    return false
end

function assertions:Failf(t, fmt, ...)
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
    return self:Fail(t, string.format([[expected "%s"; got "%s"]] .. '\n', expected, actual), ...)
end

function assertions:Equalf(t, expected, actual, fmt, ...)
    if expected == actual then
        return true
    end
    return self:Failf(t, string.format([[expected "%s"; got "%s"%s%s]], expected, actual, '\n', fmt), ...)
end

function assertions:NotEqual(t, expected, actual, ...)
    if expected ~= actual then
        return true
    end
    return self:Fail(t, string.format([[expected ~= "%s";]] .. '\n', expected), ...)
end

function assertions:NotEqualf(t, expected, actual, fmt, ...)
    if expected ~= actual then
        return true
    end
    return self:Failf(t, string.format([[expected ~= "%s"; got "%s"%s%s]], expected, '\n', fmt), ...)
end

function assertions:True(t, actual, ...)
    if actual then
        return true
    end
    return self:Fail(t, string.format([[expected true; got %s]] .. '\n', actual), ...)
end

function assertions:Truef(t, actual, fmt, ...)
    if actual then
        return true
    end
    return self:Failf(t, string.format([[expected true; got %s%s%s]], actual, '\n', fmt), ...)
end

return assertions
