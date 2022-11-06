local strings = require 'strings'

local suite = {
    t = nil,
    Suite = {},
}

function suite.Run(t, testSuite)
    -- testSuite must be subclass of suite.Suite, so will have this method
    testSuite:SetT(t)

    pcall(testSuite.SetupSuite, testSuite)
    for k, v in pairs(testSuite) do
        if strings.has_prefix(k, "Test") then
            pcall(testSuite.SetupTest, testSuite)
            t:Run(k, function(t)
                v(testSuite)
            end)
            pcall(testSuite.TearDownTest, testSuite)
        end
    end
    pcall(testSuite.TearDownSuite, testSuite)
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

return suite
