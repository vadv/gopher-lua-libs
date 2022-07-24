local runtime = require("runtime")

function Test_runtime(t)
    t:Logf("goos=%s", runtime.goos())
    t:Logf("goarch=%s", runtime.goarch())
end
