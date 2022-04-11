package io

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	lua "github.com/yuin/gopher-lua"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestWrite(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	require.NoError(t, L.DoString(`io = require("io"); return io.open("/tmp/tst.out", "w")`))
	writer := L.CheckAny(1)
	L.Pop(L.GetTop())
	wrapper := NewLuaIOWrapper(L, writer)
	//goland:noinspection GoUnhandledErrorResult
	defer wrapper.Close()
	_, err := wrapper.Write([]byte("foo bar baz"))
	require.NoError(t, err)
	wrapper.Close()
	data, err := ioutil.ReadFile("/tmp/tst.out")
	require.NoError(t, err)
	assert.Equal(t, "foo bar baz", string(data))
}

func TestRead(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	require.NoError(t, ioutil.WriteFile("/tmp/tst.in", []byte("abc def ghi"), os.ModePerm))
	require.NoError(t, L.DoString(`io = require("io"); return io.open("/tmp/tst.in", "r")`))
	reader := L.CheckAny(1)
	L.Pop(L.GetTop())
	wrapper := NewLuaIOWrapper(L, reader)
	//goland:noinspection GoUnhandledErrorResult
	defer wrapper.Close()
	data, err := ioutil.ReadAll(wrapper)
	require.NoError(t, err)
	assert.Equal(t, "abc def ghi", string(data))
}

func TestSeek(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	require.NoError(t, ioutil.WriteFile("/tmp/tst.in", []byte("abc def ghi"), os.ModePerm))
	require.NoError(t, L.DoString(`io = require("io"); return io.open("/tmp/tst.in", "r")`))
	reader := L.CheckAny(1)
	L.Pop(L.GetTop())
	wrapper := NewLuaIOWrapper(L, reader)
	//goland:noinspection GoUnhandledErrorResult
	defer wrapper.Close()
	_, err := wrapper.Seek(4, io.SeekStart)
	require.NoError(t, err)
	three := make([]byte, 3)
	num, err := wrapper.Read(three)
	require.NoError(t, err)
	require.EqualValues(t, 3, num)
	assert.Equal(t, "def", string(three))

	_, err = wrapper.Seek(1, io.SeekCurrent)
	require.NoError(t, err)
	num, err = wrapper.Read(three)
	require.NoError(t, err)
	require.EqualValues(t, 3, num)
	assert.Equal(t, "ghi", string(three))

	_, err = wrapper.Seek(0, io.SeekStart)
	require.NoError(t, err)
	num, err = wrapper.Read(three)
	require.NoError(t, err)
	require.EqualValues(t, 3, num)
	assert.Equal(t, "abc", string(three))
}
