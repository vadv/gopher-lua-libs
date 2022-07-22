# inspect [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/inspect?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/inspect)

## Usage

```lua
local optparse = require("optparse")
local opt = optparse.OptionParser { usage = "%prog [options] [gzip-file...]",
                                    version = "foo 1.23", add_help_option = false }
opt.add_option { "-h", "--help", action = "store_true", dest = "help",
                 help = "give this help" }
opt.add_option {
    "-f", "--force", dest = "force", action = "store_true",
    help = "force overwrite of output file" }

local options, args = opt.parse_args()
```
