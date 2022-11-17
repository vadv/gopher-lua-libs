local log = require 'log'

local log_levels = {
    levels = { DEBUG = 0, INFO = 1, WARN = 2, ERROR = 3 }
}
log_levels.defaultOutput = 'STDERR'
setmetatable(log_levels, {
    __index = log
})

for level in pairs(log_levels.levels) do
    log_levels[level] = log.new('STDERR')
    log_levels[level]:set_prefix(string.format('[%s] ', level))
    log_levels[level]:set_flags { date = true }
end

function log_levels.Level()
    return log_levels.level
end

function log_levels.SetLevel(level)
    level = string.upper(level)
    local level_value = log_levels.levels[level]
    if not level_value then
        error('Illegal level ' + level)
    end
    log_levels.level = level
    for k, v in pairs(log_levels.levels) do
        if level_value <= v then
            log_levels[k]:set_output(log_levels.defaultOutput)
        else
            log_levels[k]:set_output('/dev/null')
        end
    end
end

function log_levels.DefaultOutput()
    return log_levels.defaultOutput
end

function log_levels.SetDefaultOutput(output)
    log_levels.defaultOutput = output
end

return log_levels