local log = require 'log'

log.levels = { DEBUG = 0, INFO = 1, WARN = 2, ERROR = 3 }
log.defaultOutput = 'STDERR'

for level in pairs(log.levels) do
    log[level] = log.new('STDERR')
    log[level]:set_prefix(string.format('[%s] ', level))
    log[level]:set_flags { date = true }
end

function log.Level()
    return log.level
end

function log.SetLevel(level)
    level = string.upper(level)
    local level_value = log.levels[level]
    if not level_value then
        error('Illegal level ' + level)
    end
    log.level = level
    for k, v in pairs(log.levels) do
        if level_value <= v then
            log[k]:set_output(log.defaultOutput)
        else
            log[k]:set_output('/dev/null')
        end
    end
end

function log.DefaultOutput()
    return log.defaultOutput
end

function log.SetDefaultOutput(output)
    log.defaultOutput = output
end

return log