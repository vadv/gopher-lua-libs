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
