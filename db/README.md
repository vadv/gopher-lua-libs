# db [![GoDoc](https://godoc.org/github.com/vadv/gopher-lua-libs/db?status.svg)](https://godoc.org/github.com/vadv/gopher-lua-libs/db)

## Usage

```lua
local db = require("db")

local config = {
  shared = true, -- share connections between lua states
  max_connections = 1, -- max connection (if you open shared connection with different max_connections - first win)
  read_only = false,   -- must execute read-write query
}

local sqlite, err = db.open("sqlite3", "file:test.db?cache=shared&mode=memory", config)
if err then error(err) end

local result, err = sqlite:query("select 1")
if err then error(err) end
if not(result.rows[1][1] == 1) then error("sqlite error") end

local _, err = sqlite:exec("CREATE TABLE t (id int, name string);")
if err then error(err) end

for i = 1, 10 do
    local query = "INSERT INTO t VALUES ("..i..", \"name-"..i.."\");"
    if i % 2 == 0 then query = "INSERT INTO t VALUES ("..i..", NULL);" end
    local _, err = sqlite:exec(query)
    if err then error(err) end
end

local result, err = sqlite:query("select * from t;")
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

local _, err = sqlite:exec("CREATE TABLE t_stmt (id int, name string);")
if err then error(err) end

-- stmt exec
local stmt, err = sqlite:stmt("insert into t_stmt (id, name) values (?, ?)")
if err then error(err) end
local result, err = stmt:exec(1, 'name-1')
if err then error(err) end
if not(result.rows_affected == 1) then error("affted: "..tostring(result.rows_affected)) end
local err = stmt:close()
if err then error(err) end

-- stmt query
local stmt, err = sqlite:stmt("select name from t_stmt where id = ?")
if err then error(err) end
local result, err = stmt:query(1)
if err then error(err) end
if not(result.rows[1][1] == 'name-1') then error("must be 'name-1': "..tostring(result.rows[1][1])) end
local err = stmt:close()
if err then error(err) end

-- command (outside transaction)
local _, err = sqlite:command("PRAGMA journal_mode = OFF;")
if err then error(err) end

local err = sqlite:close()
if err then error(err) end
```

## Supported Drivers

* [sqlite3](https://github.com/mattn/go-sqlite3)
* [postgres](https://github.com/lib/pq/)
* [mysql](https://github.com/go-sql-driver/mysql)
