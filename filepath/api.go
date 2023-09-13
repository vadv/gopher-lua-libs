// Package filepath implements golang filepath functionality for lua.
package filepath

import (
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

// Abs returns an absolute representation of path.
func Abs(L *lua.LState) int {
	path := L.CheckString(1)
	ret, err := filepath.Abs(path)
	L.Push(lua.LString(ret))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// Basename lua filepath.basename(path) returns the last element of path
func Basename(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.Base(path)))
	return 1
}

// Clean returns the shortest path name equivalent to path
func Clean(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.Clean(path)))
	return 1
}

// Dir lua filepath.dir(path) returns all but the last element of path, typically the path's directory
func Dir(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.Dir(path)))
	return 1
}

// EvalSymlinks returns the path name after the evaluation of any symbolic link.
func EvalSymlinks(L *lua.LState) int {
	path := L.CheckString(1)
	ret, err := filepath.EvalSymlinks(path)
	L.Push(lua.LString(ret))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// FromSlash returns the result of replacing each slash ('/') character
// in path with a separator character. Multiple slashes are replaced
// by multiple separators.
func FromSlash(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.FromSlash(path)))
	return 1
}

// Ext lua filepath.ext(path) returns the file name extension used by path.
func Ext(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.Ext(path)))
	return 1
}

// Glob: filepath.glob(pattern) returns the names of all files matching pattern or nil if there is no matching file.
func Glob(L *lua.LState) int {
	pattern := L.CheckString(1)
	files, err := filepath.Glob(pattern)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	result := L.CreateTable(len(files), 0)
	for _, file := range files {
		result.Append(lua.LString(file))
	}
	L.Push(result)
	return 1
}

// IsAbs reports whether the path is absolute.
func IsAbs(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LBool(filepath.IsAbs(path)))
	return 1
}

// Join lua fileapth.join(path, ...) joins any number of path elements into a single path, adding a Separator if necessary.
func Join(L *lua.LState) int {
	var elems []string
	for i := 1; i <= L.GetTop(); i++ {
		elem := L.CheckAny(i).String()
		elems = append(elems, elem)
	}
	path := filepath.Join(elems...)
	L.Push(lua.LString(path))
	return 1
}

// ListSeparator lua filepath.list_separator() OS-specific path list separator
func ListSeparator(L *lua.LState) int {
	L.Push(lua.LString(filepath.ListSeparator))
	return 1
}

// Match reports whether name matches the shell file name pattern.
func Match(L *lua.LState) int {
	pattern := L.CheckString(1)
	name := L.CheckString(2)
	matched, err := filepath.Match(pattern, name)
	L.Push(lua.LBool(matched))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// Rel returns a relative path
func Rel(L *lua.LState) int {
	basepath := L.CheckString(1)
	targpath := L.CheckString(2)
	ret, err := filepath.Rel(basepath, targpath)
	L.Push(lua.LString(ret))
	if err != nil {
		L.Push(lua.LString(err.Error()))
		return 2
	}
	return 1
}

// Separator lua filepath.separator() OS-specific path separator
func Separator(L *lua.LState) int {
	L.Push(lua.LString(filepath.Separator))
	return 1
}

// Split splits path immediately following the final Separator,
// separating it into a directory and file name component.
func Split(L *lua.LState) int {
	path := L.CheckString(1)
	dir, file := filepath.Split(path)
	L.Push(lua.LString(dir))
	L.Push(lua.LString(file))
	return 2
}

func SplitList(L *lua.LState) int {
	path := L.CheckString(1)
	ret := filepath.SplitList(path)
	table := L.CreateTable(len(ret), 0)
	for _, part := range ret {
		table.Append(lua.LString(part))
	}
	L.Push(table)
	return 1
}

// ToSlash returns the result of replacing each separator character
// in path with a slash ('/') character. Multiple separators are
// replaced by multiple slashes.
func ToSlash(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.ToSlash(path)))
	return 1
}

// VolumeName returns leading volume name.
// Given "C:\foo\bar" it returns "C:" on Windows.
// Given "\\host\share\foo" it returns "\\host\share".
// On other platforms it returns "".
func VolumeName(L *lua.LState) int {
	path := L.CheckString(1)
	L.Push(lua.LString(filepath.VolumeName(path)))
	return 1
}
