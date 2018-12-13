package tac

import (
	"bufio"
	"errors"
	"io"
)

const maxBufSize = 8 * 1024

type tacScanner struct {
	r            io.ReadSeeker
	splitFunc    bufio.SplitFunc
	buf          []byte
	offset       int64
	atEOF        bool
	tokens       [][]byte
	partialToken int
	lastErr      error
}

func newTacScanner(r io.ReadSeeker) *tacScanner {
	b := &tacScanner{
		r:         r,
		buf:       make([]byte, 4096),
		atEOF:     true,
		splitFunc: bufio.ScanLines,
	}
	b.offset, b.lastErr = r.Seek(0, 2)
	return b
}

func (b *tacScanner) fillbuf() error {
	b.tokens = b.tokens[:0]
	if b.offset == 0 {
		return io.EOF
	}
	space := len(b.buf) - b.partialToken
	if space == 0 {
		if len(b.buf) >= maxBufSize {
			return errors.New("token too long")
		}
		n := len(b.buf) * 2
		if n > maxBufSize {
			n = maxBufSize
		}
		newBuf := make([]byte, n)
		copy(newBuf, b.buf[0:b.partialToken])
		b.buf = newBuf
		space = len(b.buf) - b.partialToken
	}
	if int64(space) > b.offset {
		b.buf = b.buf[0 : b.partialToken+int(b.offset)]
		space = len(b.buf) - b.partialToken
	}
	newOffset := b.offset - int64(space)
	copy(b.buf[space:], b.buf[0:b.partialToken])
	_, err := b.r.Seek(newOffset, 0)
	if err != nil {
		return err
	}
	b.offset = newOffset
	if _, err := io.ReadFull(b.r, b.buf[0:space]); err != nil {
		return err
	}
	if b.offset > 0 {
		advance, _, err := b.splitFunc(b.buf, b.atEOF)
		if err != nil {
			return err
		}
		b.partialToken = advance
		if advance == 0 || advance == len(b.buf) {
			return b.fillbuf()
		}
	} else {
		b.partialToken = 0
	}
	for i := b.partialToken; i < len(b.buf); {
		advance, token, err := b.splitFunc(b.buf[i:], b.atEOF)
		if err != nil {
			b.tokens = b.tokens[:0]
			return err
		}
		if advance == 0 {
			break
		}
		b.tokens = append(b.tokens, token)
		i += advance
	}
	b.atEOF = false
	if len(b.tokens) == 0 {
		return b.fillbuf()
	}
	return nil
}

func (b *tacScanner) scan() bool {
	if len(b.tokens) > 0 {
		b.tokens = b.tokens[0 : len(b.tokens)-1]
	}
	if len(b.tokens) > 0 {
		return true
	}
	if b.lastErr != nil {
		return false
	}
	b.lastErr = b.fillbuf()
	return len(b.tokens) > 0
}

func (b *tacScanner) split(split bufio.SplitFunc) {
	b.splitFunc = split
}

func (b *tacScanner) bytes() []byte {
	return b.tokens[len(b.tokens)-1]
}

func (b *tacScanner) text() string {
	return string(b.bytes())
}

func (b *tacScanner) err() error {
	if len(b.tokens) > 0 {
		return nil
	}
	if b.lastErr == io.EOF {
		return nil
	}
	return b.lastErr
}
