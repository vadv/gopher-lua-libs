local strings = require 'strings'

local suite = {
    t = nil,
    Suite = {},
}

function suite.Run(t, testSuite)
    local testCount = 0

    -- testSuite must be subclass of suite.Suite, so will have this method
    testSuite:SetT(t)

    if testSuite.SetupSuite then
        testSuite:SetupSuite()
    end
    for k, v in pairs(testSuite) do
        if strings.has_prefix(k, "Test") then
            testCount = testCount + 1
            if testSuite.SetupTest then
                testSuite:SetupTest()
            end
            t:Run(k, function(t)
                testSuite:SetT(t)
                v(testSuite)
            end)
            testSuite:SetT(t)
            if testSuite.TearDownTest then
                testSuite:TearDownTest()
            end
        end
    end
    if testSuite.TearDownSuite then
        testSuite:TearDownSuite()
    end

    return testCount
end

function suite.Suite:new(o)
    o = o or {}   -- create object if user does not provide one
    setmetatable(o, self)
    self.__index = self
    return o
end

function suite.Suite:T()
    return self.t
end

function suite.Suite:SetT(t)
    self.t = t
end

function suite.Suite:Run(name, f)
    local t = self:T()
    self:T():Run(name, function(t)
        self:SetT(t)
        f(self)
    end)
    self:SetT(t)
end

return suite
