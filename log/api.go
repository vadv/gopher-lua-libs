// Package log provides logging features to lua.
package log

import (
	"log"
	"os"

	lua "github.com/yuin/gopher-lua"
)

const (
	// Default format for Flags.
	DefaultLogFlags = log.LstdFlags
	// Default minimum level for logging
	DefaultLogLevel = `warn`
)

var (
	AvaliableLevels = []string{
		`fatal`,
		`error`,
		`warn`,
		`info`,
		`debug`,
	}
)

// Loggers contains all avaliable Log.Logger instances and map with enabled loggers
type Loggers struct {
	Debug         *log.Logger
	Info          *log.Logger
	Warn          *log.Logger
	Error         *log.Logger
	Fatal         *log.Logger
	LevelsEnabled map[string]bool
}

// NewLogger(): lua logger() returns *Loggers to lua state
func NewLogger(L *lua.LState) int {
	var flags int
	var logLevel string

	if L.GetTop() > 0 {
		params := L.CheckTable(1)
		LFlags := params.RawGetString(`flags`)
		LLevel := params.RawGetString(`level`)
		if LFlags != lua.LNil {
			if value, ok := LFlags.(lua.LNumber); ok {
				flags = int(value)
			} else {
				L.ArgError(1, `type int expected for flags`)
			}
		}

		if LLevel != lua.LNil {
			if value, ok := LLevel.(lua.LString); ok {
				logLevel = value.String()
			} else {
				L.ArgError(1, `type string expected for level`)
			}
		}

	} else {
		logLevel = DefaultLogLevel
		flags = DefaultLogFlags
	}

	handlerDebug := log.New(os.Stdout, "", flags)
	handlerInfo := log.New(os.Stdout, "", flags)
	handlerWarn := log.New(os.Stderr, "", flags)
	handlerError := log.New(os.Stderr, "", flags)
	handlerFatal := log.New(os.Stderr, "", flags)

	loggersEnabled := SetLogLevel(logLevel)
	loggers := &Loggers{Info: handlerInfo,
		Warn: handlerWarn, Error: handlerError, Fatal: handlerFatal,
		Debug: handlerDebug, LevelsEnabled: loggersEnabled}

	ud := L.NewUserData()
	ud.Value = loggers
	L.SetMetatable(ud, L.GetTypeMetatable(`logger`))
	L.Push(ud)
	return 1
}

func Debug(L *lua.LState) int {
	loggers := checkLogger(L, 1)

	if checkEnabled(loggers, `debug`) == true {
		currentLogger := loggers.Debug
		msg := `[DEBUG] ` + L.CheckString(2)
		outputLog(L, currentLogger, msg)
	}

	return 1
}

func Info(L *lua.LState) int {
	loggers := checkLogger(L, 1)

	if checkEnabled(loggers, `info`) == true {
		currentLogger := loggers.Info
		msg := `[INFO] ` + L.CheckString(2)
		outputLog(L, currentLogger, msg)
	}

	return 1
}

func Warn(L *lua.LState) int {
	loggers := checkLogger(L, 1)

	if checkEnabled(loggers, `warn`) == true {
		currentLogger := loggers.Warn
		msg := `[WARN] ` + L.CheckString(2)
		outputLog(L, currentLogger, msg)
	}

	return 1
}

func Error(L *lua.LState) int {
	loggers := checkLogger(L, 1)

	if checkEnabled(loggers, `error`) == true {
		currentLogger := loggers.Error
		msg := `[ERROR] ` + L.CheckString(2)
		outputLog(L, currentLogger, msg)
	}

	return 1
}

// Unlike other levers raises an error.
func Fatal(L *lua.LState) int {
	loggers := checkLogger(L, 1)

	if checkEnabled(loggers, `fatal`) == true {
		currentLogger := loggers.Fatal
		msg := `[FATAL] ` + L.CheckString(2)
		outputLog(L, currentLogger, msg)
		L.Error(L, 1)
	}

	return 1
}

// Slice through avaliable log levels and create map with enabled
func SetLogLevel(level string) map[string]bool {
	lower := false
	LevelsEnabled := make(map[string]bool)
	for _, v := range AvaliableLevels {
		if v == level {
			lower = true
			LevelsEnabled[v] = true
			continue
		}

		if lower != true {
			LevelsEnabled[v] = true
		} else {
			LevelsEnabled[v] = false
		}
	}
	return LevelsEnabled
}

func checkLogger(L *lua.LState, n int) *Loggers {
	ud := L.CheckUserData(n)
	if v, ok := ud.Value.(*Loggers); ok {
		return v
	}
	L.ArgError(n, `type logger expected`)
	return nil
}

func checkEnabled(loggers *Loggers, level string) bool {
	return loggers.LevelsEnabled[level]
}

// Main action. Uses log.Println.
func outputLog(L *lua.LState, logger *log.Logger, msg string) {
	logger.Println(msg)
}
