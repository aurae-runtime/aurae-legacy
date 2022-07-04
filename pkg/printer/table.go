/*===========================================================================*\
 *           MIT License Copyright (c) 2022 Kris Nóva <kris@nivenly.com>     *
 *                                                                           *
 *                ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓                *
 *                ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗   ┃                *
 *                ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗  ┃                *
 *                ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║  ┃                *
 *                ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║  ┃                *
 *                ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║  ┃                *
 *                ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝  ┃                *
 *                ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛                *
 *                                                                           *
 *                       This machine kills fascists.                        *
 *                                                                           *
\*===========================================================================*/

package printer

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
)

type Table struct {
	Title         string
	i             int
	OrderedFields map[int]*Field
	width         int
	height        int
}

func NewTable(title string) *Table {
	width, height, err := terminal.GetSize(0)
	if err != nil {
		width = 80
		height = 240
	}
	return &Table{
		OrderedFields: make(map[int]*Field),
		Title:         title,
		i:             0,
		width:         width,
		height:        height,
	}
}

func (t *Table) AddField(f *Field) {
	t.OrderedFields[t.i] = f
	t.i++
}

func (t *Table) Print(w io.Writer) error {
	// Title
	fmt.Fprintf(w, "%s:\n", t.Title)

	// Print the Headers first
	headerLine := ""
	for i := 0; i <= t.i; i++ {
		f := t.OrderedFields[i]
		if f == nil {
			continue // Sanity check
		}

		headerLine += fmt.Sprintf("%*s", f.width, f.Header)
	}
	headerLine += "\n"
	fmt.Fprintf(w, headerLine)

	// Print the Fields next
	fieldLine := ""
	row := 0 // Row is the "row" of the fields to print. Start at 0
	for i := 0; i <= t.i; i++ {
		f := t.OrderedFields[i]
		if f == nil {
			continue // Sanity check
		}
		fieldLine += fmt.Sprintf("%*s", f.width, f.values[row])
		row++
	}
	fieldLine += "\n"
	fmt.Fprintf(w, fieldLine)
	return nil
}
