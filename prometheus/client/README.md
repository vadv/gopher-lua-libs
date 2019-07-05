# prometheus [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/prometheus/client?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/prometheus/client)

## Usage

```lua
local prometheus = require("prometheus")

local pp = prometheus.register(":8080")
pp:start()

local gauge = prometheus.gauge({
     namespace="node_scout",
     subsystem="nf_conntrack",
     name="insert_failed",
     help="insert_failed from nf_conntrack",
})
gauge:set(100)
```

