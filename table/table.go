package table

import (
	"strings"
<<<<<<< Updated upstream:datatypes.go
	"sync"
	"sync/atomic"
=======
	"unicode/utf8"
>>>>>>> Stashed changes:table/table.go
)

type Table struct {
	Lines   [][2]string
	Fmtfunc func(string, string, string) string
}

func New() *Table {
	return &Table{
		Lines: [][2]string{},
		Fmtfunc: func(n, v, p string) string {
			//     NAME....VALUE
			return n + p + v
		},
	}
}

func (self *Table) AddLine(name, value string) {
	self.Lines = append(self.Lines, [2]string{name, value})
}

func (self *Table) SetLine(name, value string) {
	for i := range self.Lines {
		if self.Lines[i][0] == name {
			self.Lines[i][1] = value
			return
		}
	}

	self.AddLine(name, value)
}

func (self *Table) InsertLine(line [2]string, i int) {
	if i < 0 {
		i = 0
	} else if i > len(self.Lines) {
		i = len(self.Lines)
	}

	n := make([][2]string, len(self.Lines)+1)
	n[i] = line
	copy(n[:i], self.Lines[:i])
	copy(n[i+1:len(n)], self.Lines[i:len(n)-1])

	self.Lines = n
}

func (self *Table) String() string {
	var (
		max = 0
		l   = 0
	)

	for _, x := range self.Lines {
		if len(x[0]) > max {
			max = len(x[0])
			l += len(x[1])
		}
	}

	b := strings.Builder{}
	b.Grow((len(self.Lines)*max + l) * 2)

	for _, x := range self.Lines {
		b.WriteString(self.Fmtfunc(
			x[0],
			x[1],
			strings.Repeat(" ", max-len(x[0])),
		))
		b.WriteRune('\n')
	}

	return b.String()
}
