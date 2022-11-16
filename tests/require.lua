local assertions = require 'assertions'
local orig_require = require
return assertions:new {
    fail_now = true,
    call = orig_require,
}

