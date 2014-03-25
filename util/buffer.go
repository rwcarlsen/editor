package util

import (
	"unicode/utf8"
	"bytes"
)

type Buffer struct {
	data  []byte
	lines [][]rune
}

func NewBuffer(data []byte) *Buffer {
	b := &Buffer{data: data}
	b.updLines()
	return b
}

func (b *Buffer) Rune(line, char int) rune {
	return b.lines[line][char]
}

func (b *Buffer) Line(n int) []rune {
	return b.lines[n]
}

func (b *Buffer) updLines() {
	slines := bytes.SplitAfter(b.data, []byte("\n"))
	b.lines = make([][]rune, len(slines))
	for i, l := range slines {
		b.lines[i] = bytes.Runes(l)
	}
	if len(b.lines) > 0 {
		i := len(b.lines)-1
		l := b.lines[i]
		if len(l) > 0 && l[len(l)-1] != '\n' {
			b.lines[i] = append(l, '\n')
		} else if len(b.lines[i]) == 0 {
			b.lines = b.lines[:i]
		}
	}
}

// Nlines returns the total number of lines (separated by '\n') in the buffer.
func (b *Buffer) Nlines() int {
	return len(b.lines)
}

// Insert adds passed runes into the buffer at the given byte offset.
func (b *Buffer) Insert(offset int, rs ...rune) {
	bs := []byte(string(rs))
	b.data = append(b.data[:offset], append(bs, b.data[offset:]...)...)
	b.updLines()
}

// Delete removes nrunes characters starting at the given byte offset. If
// nrunes is negative, offset is the exclusive upper bound of the removed
// characters. It returns the number of bytes removed.
func (b *Buffer) Delete(offset, nrunes int) (n int) {
	if nrunes == 0 {
		return 0
	}

	nb := 0
	if n > 0 {
		for n := 0; n < nrunes; n++ {
			_, size := utf8.DecodeRune(b.data[offset+nb:])
			nb += size
		}
		b.data = append(b.data[:offset], b.data[offset+nb:]...)
	} else {
		for n := 0; n > nrunes; n-- {
			_, size := utf8.DecodeLastRune(b.data[:offset+nb])
			nb -= size
		}
		b.data = append(b.data[:offset+nb], b.data[offset:]...)
	}
	b.updLines()
	return nb * nb / nb
}

// Pos returns the line and character index of the given byte offset.
func (b *Buffer) Pos(offset int) (line, char int) {
	lines := bytes.SplitAfter(b.data[:offset], []byte("\n"))
	return len(lines)-1, utf8.RuneCount(lines[len(lines)-1])
}

// Offset returns the byte offset of the given line and char index.
func (b *Buffer) Offset(line, char int) int {
	offset := 0
	for _, line := range b.lines[:line] {
		offset += len(line)
	}
	return offset + char
}

func (b *Buffer) Bytes() []byte {
	return b.data
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
