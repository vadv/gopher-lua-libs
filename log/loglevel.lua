local log = require 'log'

-- Set up a loglevel object that proxies to log object and adds levels and logs at each level.
local loglevel = {
    levels = { DEBUG = 0, INFO = 1, WARN = 2, ERROR = 3 },
    defaultOutput = 'STDERR',
    _level = 'INFO',
}
setmetatable(loglevel, {
    __index = log
})

-- Gets the output for the given level as compared to the loglevel._level
local function output_for_level(level)
    level = string.upper(level)
    local level_value = loglevel.levels[level]
    if not level_value then
        error('Illegal level ' + level)
    end
    local current_level_value = loglevel.levels[loglevel._level]
    local output = (current_level_value <= level_value) and loglevel.defaultOutput or '/dev/null'
    return output
end

-- Attach the logs to the loglevel object
for level in pairs(loglevel.levels) do
    local output = output_for_level(level)
    loglevel[level] = log.new(output)
    loglevel[level]:set_prefix(string.format('[%s] ', level))
    loglevel[level]:set_flags { date = true }
end

-- Returns the current level
function loglevel.get_level()
    return loglevel._level
end

-- Sets the level and adjusts all logs to either squelch or go to the defaultOutput
function loglevel.set_level(new_level)
    new_level = string.upper(new_level)
    local new_level_value = loglevel.levels[new_level]
    if not new_level_value then
        error('Illegal level ' + new_level)
    end
    loglevel._level = new_level
    for level in pairs(loglevel.levels) do
        local output = output_for_level(level)
        loglevel[level]:set_output(output)
    end
end

-- Gets the default output
function loglevel.default_output()
    return loglevel.defaultOutput
end

-- Set the default output
function loglevel.set_default_output(output)
    loglevel.defaultOutput = output
end

return loglevel