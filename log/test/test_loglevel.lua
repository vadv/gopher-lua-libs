local log = require 'loglevel'
local suite = require 'suite'
local filepath = require 'filepath'
local ioutil = require 'ioutil'
local assert = require 'assert'
local strings = require 'strings'

local LogLevelSuite = suite.Suite:new {
    stderr = io.stderr,
}

function LogLevelSuite:SetupTest()
    self.temp_dir = self:T():TempDir()
    self.output = filepath.join(self.temp_dir, 'test.output')
    log.set_default_output(self.output)
    log.set_level('INFO')
end

function LogLevelSuite:TearDownTest()
    log.set_default_output('STDERR')
    log.set_level('INFO')
end

function LogLevelSuite:getOutput()
    return ioutil.read_file(self.output)
end

function TestLogLevelSuite(t)
    assert:Greater(t, suite.Run(t, LogLevelSuite), 0, 'no tests in suite')
end

function LogLevelSuite:TestLogObjectsExist()
    assert:NotNil(self:T(), log.DEBUG)
    assert:NotNil(self:T(), log.INFO)
    assert:NotNil(self:T(), log.WARN)
    assert:NotNil(self:T(), log.ERROR)
end

function LogLevelSuite:TestDebugNoContent()
    log.DEBUG:print('foobar')
    local got, err = self:getOutput()
    assert:NoError(self:T(), err)
    assert:NotNil(self:T(), got)
    assert:Equal(self:T(), "", strings.trim_space(got))
end

function LogLevelSuite:TestDebugWithDebugSetHasContent()
    log.set_level('DEBUG')
    log.DEBUG:print('foobar')
    local got, err = self:getOutput()
    assert:NoError(self:T(), err)
    assert:NotNil(self:T(), got)
    assert:NotEqual(self:T(), "", strings.trim_space(got))
end

function LogLevelSuite:TestInfoHasContent()
    log.set_level('INFO')
    log.INFO:print('foobar')
    local got, err = self:getOutput()
    assert:NoError(self:T(), err)
    assert:NotNil(self:T(), got)
    assert:NotEqual(self:T(), "", strings.trim_space(got))
end

function LogLevelSuite:TestErrorHasContent()
    log.ERROR:print('foobar')
    local got, err = self:getOutput()
    assert:NoError(self:T(), err)
    assert:NotNil(self:T(), got)
    assert:NotEqual(self:T(), "", strings.trim_space(got))
end

function LogLevelSuite:TestBogusLogLevelHasError()
    local ok, err = pcall(log.set_level, 'DJFDJFDJFJF')
    assert:False(self:T(), ok)
    assert:Error(self:T(), err)
end

function LogLevelSuite:TestLogNew()
    -- Test that the metadata mechanism works - i.e. that the loglevel object returned can still call log methods.
    local l = log.new('abc')
    l:set_output(self.output)
    l:print('def')
    local got, err = self:getOutput()
    assert:NoError(self:T(), err)
    assert:Equal(self:T(), 'def\n', got)
end