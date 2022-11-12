local suite = require 'suite'

local MySuite = suite.Suite:new {
    setupCount = 0,
    setupSuiteCount = 0,
    tearDownCount = 0,
    tearDownSuiteCount = 0,
}

function MySuite:SetupSuite()
    self.setupSuiteCount = self.setupSuiteCount + 1
end

function MySuite:TearDownSuite()
    self.tearDownSuiteCount = self.tearDownSuiteCount + 1
end

function MySuite:SetupTest()
    self.setupCount = self.setupCount + 1
end

function MySuite:TearDownTest()
    self.tearDownCount = self.tearDownCount + 1
end

function MySuite:TestFoobar()
    -- T is available from superclass for suites; not passed in as arg
    self:T():Log('TestFoobar')
end

function MySuite:TestBaz()
    self:T():Log('TestBaz')

    self:Run('sub1', function()
        self:T():Log('sub1')
    end)

    self:Run('sub2', function()
        self:T():Log('sub2')
    end)
end

function TestSuite(t)
    -- Same mechanism for test discovery is used, but then the suite is run as sub tests via suite.Run
    assert(suite.Run(t, MySuite), "No tests were run by this Suite")

    -- Called for every test: two tests so should be 2
    assert(MySuite.setupCount == 2, tostring(MySuite.setupCount))
    assert(MySuite.tearDownCount == 2, tostring(MySuite.tearDownCount))

    -- Called only once for the suite so should be 1
    assert(MySuite.setupSuiteCount == 1, tostring(MySuite.setupSuiteCount))
    assert(MySuite.tearDownSuiteCount == 1, tostring(MySuite.tearDownSuiteCount))
end

local EmptySuite = suite.Suite:new()

function TestEmptySuiteReturnsZero(t)
    local testCount = suite.Run(t, EmptySuite)
    assert(testCount == 0, string.format('EmptySuite ran %d tests, unexpectedly', testCount))
end