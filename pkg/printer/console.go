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
	"io"
	"os"
	"strings"
)

// Console will print to the terminal using a configured output stream.
// By default, stdout will be set.
type Console struct {

	// Title is the title of this output printer
	Title string

	// OrderedPrinters will be ordered and printed
	// in order.
	OrderedPrinters map[int]Printer

	i int
}

func NewConsole(title string, a ...interface{}) *Console {
	title = fmt.Sprintf(title, a...)
	return &Console{
		i:               0,
		Title:           title,
		OrderedPrinters: make(map[int]Printer),
	}
}

func (c *Console) PrintStdout() error {
	return c.Print(os.Stdout)
}

func (c *Console) PrintStderr() error {
	return c.Print(os.Stdout)
}

func (c *Console) Print(w io.Writer) error {

	fmt.Fprintf(w, drawLine("━"))

	// Title
	if c.Title != "" {
		col := color.New(color.Bold, color.FgGreen)
		fmt.Fprintf(w, "[%s]\n", col.Sprintf(strings.ToUpper(c.Title)))
	}

	// Printers
	fmt.Fprintf(w, drawLine("─"))
	for i := 0; i <= c.i; i++ {
		p := c.OrderedPrinters[i]
		if p == nil {
			continue
		}
		err := p.Print(w)
		fmt.Fprintf(w, drawLine("─"))

		if err != nil {
			return err
		}
	}

	//fmt.Fprintf(w, drawLine("━"))

	return nil
}

func (c *Console) AddPrinter(printer Printer) {
	c.OrderedPrinters[c.i] = printer
	c.i++
}

type Printer interface {
	Print(w io.Writer) error
}
