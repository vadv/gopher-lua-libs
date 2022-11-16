local assertions = require 'assertions'
local orig_assert = assert
return assertions:new {
    call = orig_assert,
}
