package pb

import (
	"log"

	time "github.com/vadv/gopher-lua-libs/time"
	lua "github.com/yuin/gopher-lua"
)

func ExampleAllParams() {
	state := lua.NewState()
	Preload(state)
	time.Preload(state)
	source := `
                local pb = require('pb')
                local time = require('time')

                local count = 2
                local bar = pb.new(count)
                local template = string.format('%s {{ counters . }} {{percent . }} {{ etime . }}', '[custom template]')

                err = bar:configure({writer='stdout', refresh_rate=2001, template=template})
                if err then error(err) end
                bar:start()

                for i=1, count, 1 do
                  time.sleep(1)
                  bar:increment()
                end
                bar:finish()
        `
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// [custom template] 2 / 2 100.00% 2s

}
