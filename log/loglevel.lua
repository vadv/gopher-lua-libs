local log = require 'log'

local loglevel = {
    levels = { DEBUG = 0, INFO = 1, WARN = 2, ERROR = 3 }
}
loglevel.defaultOutput = 'STDERR'
setmetatable(loglevel, {
    __index = log
})

for level in pairs(loglevel.levels) do
    loglevel[level] = log.new('STDERR')
    loglevel[level]:set_prefix(string.format('[%s] ', level))
    loglevel[level]:set_flags { date = true }
end

function loglevel.Level()
    return loglevel.level
end

function loglevel.SetLevel(level)
    level = string.upper(level)
    local level_value = loglevel.levels[level]
    if not level_value then
        error('Illegal level ' + level)
    end
    loglevel.level = level
    for k, v in pairs(loglevel.levels) do
        if level_value <= v then
            loglevel[k]:set_output(loglevel.defaultOutput)
        else
            loglevel[k]:set_output('/dev/null')
        end
    end
end

function loglevel.DefaultOutput()
    return loglevel.defaultOutput
end

function loglevel.SetDefaultOutput(output)
    loglevel.defaultOutput = output
end

return loglevel