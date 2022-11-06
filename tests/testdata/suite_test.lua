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
end

function TestSuite(t)
    -- Same mechanism for test discovery is used, but then the suite is run as sub tests via suite.Run
    suite.Run(t, MySuite)

    -- Called for every test: two tests so should be 2
    assert(MySuite.setupCount == 2, tostring(MySuite.setupCount))
    assert(MySuite.tearDownCount == 2, tostring(MySuite.tearDownCount))

    -- Called only once for the suite so should be 1
    assert(MySuite.setupSuiteCount == 1, tostring(MySuite.setupSuiteCount))
    assert(MySuite.tearDownSuiteCount == 1, tostring(MySuite.tearDownSuiteCount))
end
