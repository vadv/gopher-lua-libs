local plugin = require 'plugin'
local time = require 'time'
local ioutil = require 'ioutil'
local inspect = require 'inspect'

function Test_plugin(t)
    t:Run("no payload", function(t)
        local plugin_body_1 = [[

local http = require("http")
local time = require("time")

local i = 1
local client = http.client()
local request = http.request("GET", "http://google.com")
os.remove("./test/file.txt")

while true do
    local result, err = client:do_request(request)
    if err then error(err) end
    time.sleep(0.1)
    local file = io.open("./test/file.txt", "w")
    file:write(tostring(i), "\n")
    file:close()
    i = i + 1
end

]]

        local curl_plugin = plugin.do_string(plugin_body_1)

        curl_plugin:run()
        time.sleep(2)
        curl_plugin:stop()
        time.sleep(1)

        local running = curl_plugin:is_running()
        assert(not running, "already running")

        local data, err = ioutil.read_file("./test/file.txt")
        assert(not err, err)

        t:Log(data)
        local i = tonumber(data)
        assert(i >= 1, string.format("%d < 1", i))
    end)

    t:Run("with payload", function(t)
        -- test with payload
        local plugin_body_2 = [[
    local ioutil = require("ioutil")
    local err = ioutil.write_file("./test/payload.txt", payload)
    if err then error(err) end
]]
        local test_payload = "OK"
        local plugin_with_payload = plugin.do_string_with_payload(plugin_body_2, test_payload)
        plugin_with_payload:run()

        time.sleep(1)
        local i = 0
        while plugin_with_payload:is_running() do
            time.sleep(1)
            i = i + 1
            assert(i >= 3, "timeout: i=" .. i)
        end

        local data, err = ioutil.read_file("./test/payload.txt")
        assert(not err, err)
        assert(data == test_payload, "data <> test_payload")
    end)
end

function TestArgs(t)
    local plugin_body = [[
        local ch, msg = unpack(arg)
        assert(ch, tostring(ch))
        assert(msg, tostring(msg))
        ch:send(msg.." pong")
        ch:close()
    ]]
    local ch = channel.make(1)
    local args_plugin = plugin.do_string(plugin_body, ch, "ping")
    args_plugin:run()
    time.sleep(0.1)
    local err = args_plugin:error()
    assert(not err, err)
    local ok, answer = ch:receive()
    assert(ok)
    assert(answer == "ping pong", answer)
    args_plugin:stop()
end

function TestWait(t)
    local plugin_quick_body = [[
        local time = require 'time'
        time.sleep(0.1)
    ]]
    local plugin_quick_body_fail = [[
        error('fail')
    ]]
    local plugin_slow_body = [[
        local time = require 'time'
        time.sleep(5)
    ]]

    t:Run("no timeout", function(t)
        local notimeout_plugin = plugin.do_string(plugin_quick_body)
        local err = notimeout_plugin:run()
        assert(not err, err)
        err = notimeout_plugin:wait()
        assert(not err, err)
    end)

    t:Run("no timeout fails", function(t)
        local notimeout_plugin = plugin.do_string(plugin_quick_body_fail)
        local err = notimeout_plugin:run()
        assert(not err, err)
        err = notimeout_plugin:wait()
        assert(err)
    end)

    t:Run("timeout ok", function(t)
        local notimeout_plugin = plugin.do_string(plugin_quick_body)
        local err = notimeout_plugin:run()
        assert(not err, err)
        err = notimeout_plugin:wait(1)
        assert(not err, err)
    end)

    t:Run("timeout expires", function(t)
        local notimeout_plugin = plugin.do_string(plugin_slow_body)
        local err =notimeout_plugin:run()
        assert(not err, err)
        err = notimeout_plugin:wait(0.1)
        assert(err)
    end)
end

function TestMultipleWorkers(t)
    -- Fire up 5 work consumers that double each unit of work and return {work, work * 2}
    local work_body = [[
        local workCh, resultCh = unpack(arg)

        local ok, work = workCh:receive()
        while ok do
            resultCh:send {work, work * 2}
            ok, work = workCh:receive()
        end
    ]]
    local workers = {}
    local worker_channels = {}
    local workCh = channel.make(100)
    local resultCh = channel.make(100)
    for i = 1, 5 do
        worker_plugin = plugin.do_string(work_body, workCh, resultCh)
        local err = worker_plugin:run()
        assert(not err, err)
        table.insert(workers, worker_plugin)
        table.insert(worker_channels, worker_plugin:done_channel())
    end

    -- Fire up a watcher to close the resultCh when all workers have exited
    local worker_watcher_body = [[
        resultCh, worker_channels = unpack(arg)
        for _, ch in ipairs(worker_channels) do
            ch:receive()
        end
        resultCh:close()
    ]]
    local worker_watcher_plugin = plugin.do_string(worker_watcher_body, resultCh, worker_channels)
    err = worker_watcher_plugin:run()
    assert(not err, err)

    -- Fire up a producer of work
    local work_producer_body = [[
        local workCh = unpack(arg)
        for i = 1, 10 do
            workCh:send(i)
        end
        workCh:close()
    ]]
    local work_producer_plugin = plugin.do_string(work_producer_body, workCh)
    err = work_producer_plugin:run()
    assert(not err, err)

    -- Now just walk the results, which should close when all workers have exited
    local count = 0
    local ok, result = resultCh:receive()
    while ok do
        assert(ok)
        assert(result[1] * 2 == result[2], inspect(result))
        count = count + 1
        t:Log(inspect(result))
        ok, result = resultCh:receive()
    end
    t:Logf("count = %d", count)
    assert(count == 10, tostring(count))

    -- Ensure that all workers are terminated at this point and with no errors
    assert(not worker_watcher_plugin:is_running(), "worker_watcher_plugin is still running")
    assert(not worker_watcher_plugin:error(), worker_watcher_plugin:error())
    assert(not work_producer_plugin:is_running(), "worker_watcher_plugin is still running")
    assert(not work_producer_plugin:error(), work_producer_plugin:error())
    for i, worker in ipairs(workers) do
        assert(not worker:is_running(), string.format("worker %d is running", i))
        assert(not worker:error(), string.format("worker %d error %s", i, worker:error()))
    end
end
