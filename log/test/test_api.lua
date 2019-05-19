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

logger = log.logger({level='error'})
logger:warn('3rd WARN') -- never should be printed
logger:error('3rd ERROR')
