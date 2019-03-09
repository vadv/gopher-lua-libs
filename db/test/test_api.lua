local db = require("db")
local time = require("time")
local inspect = require("inspect")

local sqlite, err = db.open("sqlite3", "file:testdb.db?mode=memory", {shared = true} )
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

local _, err = sqlite:exec("CREATE TABLE table_time (id int, time DATETIME DEFAULT CURRENT_TIMESTAMP);")
if err then error(err) end

for i = 1, 10 do
    local query = "INSERT INTO table_time VALUES ("..i..", " .. time.unix() .. ");"
    if i % 2 == 0 then
        query = "INSERT INTO table_time(id) VALUES ("..i..");"
    end
    local result, err = sqlite:exec(query)
    if err then error(err) end
    print(inspect(result))
end

local result, err = sqlite:query("select * from table_time;")
if err then error(err) end

for _, row in pairs(result.rows) do
    for id, name in pairs(result.columns) do
        local datetime = os.date("*t", row[id])
        print(name, datetime.year, datetime.month, datetime.day, datetime.hour, datetime.sec)
    end
end

local _, err = sqlite:exec("CREATE TABLE t_stmt (id int, name string);")
if err then error(err) end

local stmt, err = sqlite:stmt("insert into t_stmt (id, name) values (?, ?)")
if err then error(err) end
local result, err = stmt:exec(1, 'name-1')
if err then error(err) end
if not(result.rows_affected == 1) then error("affted: "..tostring(result.rows_affected)) end
local err = stmt:close()
if err then error(err) end

local result, err = sqlite:query("select name from t_stmt where id = 1")
if err then error(err) end
if not(result.rows[1][1] == 'name-1') then error("must be 'name-1': "..tostring(result.rows[1][1])) end

local stmt, err = sqlite:stmt("select name from t_stmt where id = ?")
if err then error(err) end
local result, err = stmt:query(1)
if err then error(err) end
if not(result.rows[1][1] == 'name-1') then error("must be 'name-1': "..tostring(result.rows[1][1])) end
local err = stmt:close()
if err then error(err) end

-- test shared connections
local sqliteShared, err = db.open("sqlite3", "file:testdb.db?mode=memory", {shared = true})
if err then error(err) end
local result, err = sqliteShared:query("select name from t_stmt where id = 1")
if err then error(err) end
if not(result.rows[1][1] == 'name-1') then error("must be 'name-1': "..tostring(result.rows[1][1])) end
local sqliteShared2, err = db.open("sqlite3", "file:testdb.db?mode=memory", {shared = false})
if err then error(err) end
local result, err = sqliteShared2:query("select name from t_stmt where id = 1")
if not err then error("must be unknown table") end

local err = sqlite:close()
if err then error(err) end
local result, err = sqliteShared:query("select name from t_stmt where id = 1")
if not err then error("must be closed") end

