package util

import (
	"strings"
)

type Buffer struct {
	data  []rune
	lines [][]rune
}

func NewBuffer(data []byte) *Buffer {
	s := string(data)
	b := &Buffer{data: []rune(s)}
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
	slines := strings.SplitAfter(string(b.data), "\n")
	b.lines = make([][]rune, len(slines))
	for i, l := range slines {
		b.lines[i] = []rune(l)
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

func (b *Buffer) Nlines() int {
	return len(b.lines)
}

func (b *Buffer) Insert(offset int, rs ...rune) {
	b.data = append(b.data[:offset], append(rs, b.data[offset:]...)...)
	b.updLines()
}

func (b *Buffer) Delete(start, end int) {
	b.data = append(b.data[:start], b.data[end:]...)
	b.updLines()
}

func (b *Buffer) Pos(offset int) (line, char int) {
	for _, r := range b.data[:offset] {
		char++
		if r == '\n' {
			line++
			char = 0
		}
	}
	return line, char
}

func (b *Buffer) Offset(line, char int) int {
	offset := 0
	for _, line := range b.lines[:line] {
		offset += len(line)
	}
	return offset + char
}

func (b *Buffer) Bytes() []byte {
	return []byte(string(b.data))
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
