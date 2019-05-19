## Functions

 - `logger([flags], [level])` - return logger with given format and level. `level(string)` defines minimail level to handle message. `flags(int)` defines (log Flags)[https://golang.org/pkg/log/#pkg-constants] via bitmask.

## Methods
 - `debug`
    The DEBUG level designates fine-grained informational events that are most useful to debug an application. Sends messages to stdout.
 - `info`
    The INFO level designates informational messages that highlight the progress of the application at coarse-grained level.
    Sends messages to stdout
 - `warn`
    The WARN level designates potentially harmful situations.
    Sends messages to stderr
 - `error`
    The ERROR level designates error events that might still allow the application to continue running.
    Sends messages to stderr
 - `fatal`
    The FATAL level designates very severe error events that would presumably lead the application to abort. Sends messages to stderr and raise lua Error



## Examples

```lua
local log = require("log")

logger = log.logger() -- default value

logger:debug('DEBUG')
logger:info('INFO')
logger:warn('WARN')
logger:error('ERROR')
--logger:fatal('FATAL')


--2019/05/19 20:15:03 [WARN] WARN
--2019/05/19 20:15:03 [ERROR] ERROR
--2019/05/19 20:15:03 [FATAL] FATAL
-- stack traceback:
--   	[G]: in function 'fatal'
--     	./test/test_api.lua:20: in main chunk
--     	[G]: ?


-- Change format and set minimal level
logger = log.logger({flags=100000, level='debug'})
logger:debug('2nd DEBUG')
logger:info('2nd INFO')
logger:warn('2nd WARN')
logger:error('2nd ERROR')

--[DEBUG] DEBUG
--[INFO] INFO
--[WARN] WARN
--[ERROR] ERROR

```
