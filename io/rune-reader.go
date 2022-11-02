package io

import (
	"io"
	"strings"
	"unicode/utf8"
)

//UnbufferedRuneReader doesn't attempt to buffer the underlying reader, but does buffer enough to decode utf8 runes.
type UnbufferedRuneReader struct {
	reader io.Reader
	buf    [utf8.UTFMax]byte // used only inside ReadRune
	single [1]byte           // to read one byte
}

//ReadLine reads a line from any reader.
func ReadLine(reader io.Reader) (string, error) {
	rr := ToRuneReader(reader)
	var sb strings.Builder
	var r rune
	var err error
	for r, _, err = rr.ReadRune(); err == nil && r != '\n'; r, _, err = rr.ReadRune() {
		sb.WriteRune(r)
	}
	if err == io.EOF && sb.Len() > 0 {
		err = nil
	}
	return sb.String(), err
}

//ToRuneReader Converts reader to an io.RuneReader
func ToRuneReader(reader io.Reader) io.RuneReader {
	if ret, ok := reader.(io.RuneReader); ok {
		return ret
	}
	return &UnbufferedRuneReader{
		reader: reader,
	}
}

func (u *UnbufferedRuneReader) readByte() (b byte, err error) {
	n, err := io.ReadFull(u.reader, u.single[:])
	if n != 1 {
		return 0, err
	}
	return u.single[0], err
}

//ReadRune reads a single rune, and returns the rune, its byte-length, and possibly an error.
// see the code for fmt.Scanln - which is not public, but which does, but tokenizing on space, which is not desirable.
// The implementation in fmt.Scanln also implements io.RuneScanner, which is not needed here as newlines are discarded.
func (u *UnbufferedRuneReader) ReadRune() (r rune, size int, err error) {
	u.buf[0], err = u.readByte()
	if err != nil {
		return
	}
	if u.buf[0] < utf8.RuneSelf { // fast check for common ASCII case
		r = rune(u.buf[0])
		size = 1
		return
	}
	var n int
	for n = 1; !utf8.FullRune(u.buf[:n]); n++ {
		u.buf[n], err = u.readByte()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			return
		}
	}
	r, size = utf8.DecodeRune(u.buf[:n])
	return
}
