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
	"github.com/fatih/color"
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

func NewTable(title string, a ...interface{}) *Table {
	title = fmt.Sprintf(title, a...)
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
	if t.Title != "" {
		fmt.Fprintf(w, "%s\n", color.GreenString(t.Title))
	}

	// Print the Headers first
	headerLine := " " // Offset a single space
	maxRow := 0       // maxRow is dynamically calculated as we print the headers
	for i := 0; i <= t.i; i++ {
		f := t.OrderedFields[i]
		if f == nil {
			continue // Sanity check
		}
		if f.i > maxRow {
			maxRow = f.i
		}

		headerLine += fmt.Sprintf("%-*s", f.width, color.BlueString(f.Header))
	}
	headerLine += "\n"
	fmt.Fprintf(w, headerLine)

	// Print the Fields next
	// Row is the "row" of the fields to print. Start at 0
	for row := 0; row <= maxRow; row++ {
		fieldLine := ""

		for i := 0; i <= t.i; i++ {
			f := t.OrderedFields[i]
			if f == nil {
				continue // Sanity check
			}
			fieldLine += fmt.Sprintf("%*s", f.width, f.values[row])
		}
		fieldLine += "\n"
		fmt.Fprintf(w, fieldLine)
	}
	return nil
}
