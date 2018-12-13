// build +purego

package db

import (
	"log"

	inspect "github.com/vadv/gopher-lua-libs/inspect"
	lua "github.com/yuin/gopher-lua"
)

// db_ud:query()
func Example_package() {
	state := lua.NewState()
	Preload(state)
	inspect.Preload(state)
	source := `
    local db = require("db")
    local inspect = require("inspect")

    local ql, err = db.open("ql-mem", "memory://in-memory.db")
    if err then error(err) end

    local result, err = ql:query("select \"ok\";")
    if err then error(err) end
    print(inspect(result.rows))

    local _, err = ql:exec("CREATE TABLE t (id int, name string);")
    if err then error(err) end

    local result, err = ql:exec("INSERT INTO t VALUES (1, \"chook\");")
    if err then error(err) end
    print(inspect(result, {newline="", indent=""}))

    local result, err = ql:exec("INSERT INTO t VALUES (2, \"gek\");")
    if err then error(err) end
    print(inspect(result, {newline="", indent=""}))

    local result, err = ql:query("select * from t order by id desc;")
    if err then error(err) end

    print(inspect(result.columns))

    for _, row in pairs(result.rows) do
        print(inspect(row))
    end

`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// { { "ok" } }
	// {last_insert_id = 1,rows_affected = 1}
	// {last_insert_id = 2,rows_affected = 1}
	// { "id", "name" }
	// { 2, "gek" }
	// { 1, "chook" }
}
