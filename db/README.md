# db [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/db?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/db)

## Usage

```lua
local db = require("db")

local ql, err = db.open("ql-mem", "memory://mem.db")
if err then error(err) end

local result, err = ql:query("select 1")
if err then error(err) end
if not(result.rows[1][1] == 1) then error("ql error") end

local _, err = ql:exec("CREATE TABLE t (id int, name string);")
if err then error(err) end

for i = 1, 10 do
    local query = "INSERT INTO t VALUES ("..i..", \"name-"..i.."\");"
    if i % 2 == 0 then query = "INSERT INTO t VALUES ("..i..", NULL);" end
    local _, err = ql:exec(query)
    if err then error(err) end
end

local result, err = ql:query("select * from t;")
if err then error(err) end

for i, v in pairs(result.columns) do
    if i == 1 then if not(v == "id") then error("error") end end
    if i == 2 then if not(v == "name") then error("error") end end
end

for _, row in pairs(result.rows) do
    for id, name in pairs(result.columns) do
        print(name, row[id])
    end
end

local err = ql:close()
if err then error(err) end
```

## Supported Drivers

* [ql](https://godoc.org/modernc.org/ql)
* [postgres](https://github.com/lib/pq/)
* [mysql](https://github.com/go-sql-driver/mysql)
