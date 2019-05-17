// Package cmd implements golang cmd functionality for lua.
package cmd

import (
	"bytes"
	"os/exec"
	"runtime"
	"syscall"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const (
	//Timeout default execution timeout in seconds
	Timeout = 10
)

// Exec lua cmd.exec(command) return ({status=0, stdout="", stderr=""}, err)
func Exec(L *lua.LState) int {
	command := L.CheckString(1)
	timeout := time.Duration(L.OptInt64(2, Timeout)) * time.Second
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux", "darwin":
		cmd = exec.Command("sh", "-c", command)
	case "windows":
		cmd = exec.Command("cmd.exe", "/C", command)
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString(`unsupported os`))
		return 2
	}

	stdout, stderr := bytes.Buffer{}, bytes.Buffer{}
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Start(); err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		go cmd.Process.Kill()
		L.Push(lua.LNil)
		L.Push(lua.LString(`execute timeout`))
		return 2
	case err := <-done:
		result := L.NewTable()
		L.SetField(result, "stdout", lua.LString(stdout.String()))
		L.SetField(result, "stderr", lua.LString(stderr.String()))
		L.SetField(result, "status", lua.LNumber(-1))

		if err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					L.SetField(result, "status", lua.LNumber(int64(status.ExitStatus())))
				}
			}
		} else {
			L.SetField(result, "status", lua.LNumber(0))
		}
		L.Push(result)
		return 1
	}

}
