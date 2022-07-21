local pb = require('pb')
local time = require('time')

local count = 2

local function run(bar)
    bar:start()

    for i = 1, count, 1 do
        time.sleep(1)
        bar:increment()
    end

    bar:finish()
end

function Test_predefined_template(t)
    local bar = pb.new(count)
    bar:configure({ template = 'simple' })

    run(bar)
end

function Test_custom_template(t)
    local tmpl = string.format('%s {{ counters . }} {{bar . }} {{percent . }} {{ etime . }}', 'THIS IS PREFIX')

    local bar = pb.new(count)
    bar:configure({ template = tmpl })
    run(bar)
end
