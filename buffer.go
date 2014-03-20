package main

import (
)

type Buffer struct {
	data []rune
	lines = [][]rune
}

func NewBuffer(data []byte) *Buffer {
	b := &Buffer{data: []rune(data)}
	b.updLines()
	return b
}

func (b *Buffer) Rune(line, char int) rune {
	return b.lines[line][char]
}

func (b *Buffer) Line(n int) []rune {
	return b.lines[line]
}

func (b *Buffer) updLines() int {
	slines := strings.Split(string(b.data), "\n")
	b.lines := make([][]rune, len(slines))
	for i, l := range slines {
		b.lines[i] = []rune(l)
	}
}

func (b *Buffer) Nlines() int {
	return len(b.lines)
}

func (b *Buffer) Insert(offset int, rs []rune) {
	b.data = append(b.data[:offset], append(rs, b.data[offset:]...)...)
	b.updLines()
}

func (b *Buffer) Delete(start, end int) {
	b.data = append(b.data[:start], b.data[end:]...)
	b.updLines()
}

func (b *Buffer) Pos(offset int) (line, char int) {
	offset := 0
	for _, r := range b.data[:offset] {
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
		offset += len(line) + 1 // +1 for newline
	}
	offset += min(char, len(b.lines[line]))
	return offset
}

func (b *Buffer) Bytes() []byte {
	return []byte(b.data)
}

