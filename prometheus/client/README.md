# prometheus [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/prometheus/client?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/prometheus/client)

## Usage

```lua
local prometheus = require("prometheus")

local pp = prometheus.register(":8080")
pp:start()

-- gauge / counter
local gauge = prometheus.gauge({ -- prometheus.counter
     namespace="node_scout",
     subsystem="nf_conntrack",
     name="insert_failed",
     help="insert_failed from nf_conntrack",
})
gauge:set(100)
gauge:inc()
gauge:add(1)

-- gauge vector / counter vector
local gauge = prometheus.gauge({  -- prometheus.counter
     namespace="node_scout",
     subsystem="nf_conntrack",
     name="insert_failed",
     help="insert_failed from nf_conntrack",
     labels = {"label_1", "label_2"}
})
gauge:set(100, {"label_1":"one", "label_2":"two"})
gauge:inc({"label_1":"one", "label_2":"two"})
gauge:add(1, {"label_1":"one", "label_2":"two"})
```

