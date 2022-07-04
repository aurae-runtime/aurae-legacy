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
)

type KeyValueTable struct {
	Title    string
	keyWidth int
	valWidth int
	data     map[string]string
}

func NewKeyValueTable(title string) *KeyValueTable {
	return &KeyValueTable{
		Title:    title,
		data:     make(map[string]string),
		keyWidth: 0,
		valWidth: 0,
	}
}

func (t *KeyValueTable) AddKeyValue(key string, value any) {
	valStr := fmt.Sprintf("%v", value)
	if len(key) > t.keyWidth {
		t.keyWidth = len(key)
	}
	// TODO change 120 to term width
	if len(valStr) > t.valWidth && len(valStr) < 120 {
		t.valWidth = len(valStr)
	}
	t.data[key] = color.CyanString(valStr)
}

func (t *KeyValueTable) AddKeyValueErr(key string, value any) {
	valStr := fmt.Sprintf("%v", value)
	if len(key) > t.keyWidth {
		t.keyWidth = len(key)
	}
	// TODO change 120 to term width
	if len(valStr) > t.valWidth && len(valStr) < 120 {
		t.valWidth = len(valStr)
	}
	t.data[key] = color.RedString(valStr)
}

func (t *KeyValueTable) Print(w io.Writer) error {

	// Title
	if t.Title != "" {
		fmt.Fprintf(w, "%s\n", color.GreenString(t.Title))
	}

	for k, v := range t.data {
		fmt.Fprintf(w, "%-*s: %-*s", t.keyWidth*2, color.BlueString(k), t.valWidth, v)
		fmt.Fprintf(w, "\n")

	}
	return nil
}
