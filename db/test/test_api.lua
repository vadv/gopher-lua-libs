local db = require("db")
local time = require("time")
local inspect = require("inspect")

function Test_db(t)
    local sqlite, err = db.open("sqlite3", "file:testdb.db?mode=memory", { shared = true })
    assert(not err, err)

    t:Run("setup", function(t)
        local result, err = sqlite:query("select 1")
        assert(not err, err)
        assert(result.rows[1][1] == 1, result.rows[1][1])

        local _, err = sqlite:exec("CREATE TABLE t (id int, name string);")
        assert(not err, err)

        for i = 1, 10 do
            local query = "INSERT INTO t VALUES (" .. i .. ", \"name-" .. i .. "\");"
            if i % 2 == 0 then
                query = "INSERT INTO t VALUES (" .. i .. ", NULL);"
            end
            local _, err = sqlite:exec(query)
            assert(not err, err)
        end

        local result, err = sqlite:query("select * from t;")
        assert(not err, err)

        for i, v in pairs(result.columns) do
            if i == 1 then
                assert(v == "id", v .. " is not id")
            end
            if i == 2 then
                assert(v == "name", v .. " is not name")
            end
        end

        for _, row in pairs(result.rows) do
            for id, name in pairs(result.columns) do
                t:Logf("name=%s, row[%s]=%s", name, id, inspect(row[id]))
            end
        end

        local _, err = sqlite:exec("CREATE TABLE table_time (id int, time DATETIME DEFAULT CURRENT_TIMESTAMP);")
        assert(not err, err)

        for i = 1, 10 do
            local query = "INSERT INTO table_time VALUES (" .. i .. ", " .. time.unix() .. ");"
            if i % 2 == 0 then
                query = "INSERT INTO table_time(id) VALUES (" .. i .. ");"
            end
            local result, err = sqlite:exec(query)
            assert(not err, err)
            t:Log(inspect(result))
        end

        local result, err = sqlite:query("select * from table_time;")
        assert(not err, err)

        for _, row in pairs(result.rows) do
            for id, name in pairs(result.columns) do
                local datetime = os.date("*t", row[id])
                t:Log(name, datetime.year, datetime.month, datetime.day, datetime.hour, datetime.sec)
            end
        end

        local _, err = sqlite:exec("CREATE TABLE t_stmt (id int, name string);")
        assert(not err, err)

        local stmt, err = sqlite:stmt("insert into t_stmt (id, name) values (?, ?)")
        assert(not err, err)
        local result, err = stmt:exec(1, 'name-1')
        assert(not err, err)
        assert(result.rows_affected == 1, "affted: " .. tostring(result.rows_affected))
        local err = stmt:close()
        assert(not err, err)

        local result, err = sqlite:query("select name from t_stmt where id = 1")
        assert(not err, err)
        assert(result.rows[1][1] == 'name-1', "must be 'name-1': " .. tostring(result.rows[1][1]))

        local stmt, err = sqlite:stmt("select name from t_stmt where id = ?")
        assert(not err, err)
        local result, err = stmt:query(1)
        assert(not err, err)
        assert(result.rows[1][1] == 'name-1', "must be 'name-1': " .. tostring(result.rows[1][1]))
        local err = stmt:close()
        assert(not err, err)
    end)

    t:Run("shared connections", function(t)
        local sqliteShared, err = db.open("sqlite3", "file:testdb.db?mode=memory", { shared = true })
        assert(not err, err)
        local result, err = sqliteShared:query("select name from t_stmt where id = 1")
        assert(not err, err)
        assert(result.rows[1][1] == 'name-1', "must be 'name-1': " .. tostring(result.rows[1][1]))
        local sqliteShared2, err = db.open("sqlite3", "file:testdb.db?mode=memory", { shared = false })
        assert(not err, err)
        local result, err = sqliteShared2:query("select name from t_stmt where id = 1")
        assert(err, "must be unknown table")

        t:Run("command (outside transaction)", function(t)
            local _, err = sqlite:command("PRAGMA journal_mode = OFF;")
            assert(not err, err)

            local err = sqlite:close()
            assert(not err, err)
            local result, err = sqliteShared:query("select name from t_stmt where id = 1")
            assert(err, "must be closed")
        end)
    end)
end
