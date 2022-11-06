# Package tests

Support for writing lua tests like go tests; any function starting with Test will be run.

## Example go code

```go
func TestApi(t *testing.T) {
preload := tests.SeveralPreloadFuncs(
inspect.Preload,
strings.Preload,
)
assert.NotZero(t, tests.RunLuaTestFile(t, preload, "./test/test_api.lua"))
}
```

## Example lua code

```lua
function TestFoo(t)
    t:Log("foo bar baz")
    assert(someVariable, tostring(someVariable))
    expected = 2
    assert(somethingElse == expected, string.format("%d ~= %d", somethingElse, expected))
end

function TestMaybe(t)
    if os.getenv('SKIP_IT') then
        t:Skip("Skipped because SKIP_IT is defined")
    end
    assert(theActualTest, tostring(theActualTest))
end
```

## Example use of suite

To simulate [testify.suite](https://pkg.go.dev/github.com/stretchr/testify/suite) hooks, a simple suite implementation
is preloaded as well

```lua
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
```
