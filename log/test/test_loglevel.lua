local log = require 'loglevel'
local suite = require 'suite'
local filepath = require 'filepath'
local ioutil = require 'ioutil'

local LogLevelSuite = suite.Suite:new {
    stderr = io.stderr,
}

function LogLevelSuite:SetupTest()
    self.temp_dir = self:T():TempDir()
    self.output = filepath.join(self.temp_dir, 'test.output')
    log.SetDefaultOutput(self.output)
    log.SetLevel('INFO')
end

function LogLevelSuite:TearDownTest()
    log.SetDefaultOutput('STDERR')
    log.SetLevel('INFO')
end

function LogLevelSuite:getOutput()
    return ioutil.read_file(self.output)
end

function TestLogLevelSuite(t)
    assert(suite.Run(t, LogLevelSuite) > 0, 'no tests in suite')
end

function LogLevelSuite:TestLogObjectsExist()
    assert(log.DEBUG)
    assert(log.INFO)
    assert(log.WARN)
    assert(log.ERROR)
end

function LogLevelSuite:TestDebugNoContent()
    log.DEBUG:print('foobar')
    local got, err = self:getOutput()
    assert(not err, err)
    assert(got == "", string.format([[expected empty got "%s"]], got))
end

function LogLevelSuite:TestDebugWithDebugSetHasContent()
    log.SetLevel('DEBUG')
    log.DEBUG:print('foobar')
    local got, err = self:getOutput()
    assert(not err, err)
    assert(got ~= "", got)
end

function LogLevelSuite:TestInfoHasContent()
    log.SetLevel('INFO')
    log.INFO:print('foobar')
    local got, err = self:getOutput()
    assert(not err, err)
    assert(got ~= "", got)
end

function LogLevelSuite:TestErrorHasContent()
    log.ERROR:print('foobar')
    local got, err = self:getOutput()
    assert(not err, err)
    assert(got ~= "", got)
end

function LogLevelSuite:TestBogusLogLevelHasError()
    local ok, err = pcall(log.SetLevel, 'DJFDJFDJFJF')
    assert(not ok)
    assert(err)
end
