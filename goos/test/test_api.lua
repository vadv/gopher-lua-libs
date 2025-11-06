local goos = require("goos")

function Test_stat(t)
    local info, err = goos.stat("./test/test.file")
    assert(not err, err)
    assert(info.is_dir == false, "is dir")
    assert(0 == info.size, "non-zero size: " .. info.size)
    assert(info.mod_time > 0, "mod_time: " .. info.mod_time)
    assert(info.mode > "", "mode: " .. info.mode)
end

function Test_hostname(t)
    -- hostname
    local hostname, err = goos.hostname()
    assert(not err, err)
    assert(hostname > "", "hostname")
end

function Test_pagesize(t)
    assert(goos.get_pagesize() > 0, "pagesize")
end

function Test_environ(t)
    local env = goos.environ()
    assert(env, "environ should return table")
    -- Check that we get a table with environment variables
    local count = 0
    for k, v in pairs(env) do
        count = count + 1
        assert(type(k) == "string", "key should be string")
        assert(type(v) == "string", "value should be string")
    end
    assert(count > 0, "environ should return at least one environment variable")
    -- PATH should exist on most systems
    assert(env.PATH or env.Path, "PATH environment variable should exist")
    -- Test environment variable with equals sign in value
    assert(env.ENV_VAR == "TEST=1", "ENV_VAR should be TEST=1")
end
