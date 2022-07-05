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
	Title       string
	keyWidth    int
	valWidth    int
	OrderedData map[int]*Record
	i           int
}

func NewKeyValueTable(title string) *KeyValueTable {
	return &KeyValueTable{
		Title:       title,
		OrderedData: make(map[int]*Record),
		keyWidth:    0,
		valWidth:    0,
		i:           0,
	}
}

type Record struct {
	Key   string
	Value string
	Color func(format string, a ...interface{}) string
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
	t.OrderedData[t.i] = &Record{
		Key:   key,
		Value: valStr,
		Color: color.GreenString,
	}
	t.i++
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
	t.OrderedData[t.i] = &Record{
		Key:   key,
		Value: valStr,
		Color: color.RedString,
	}
	t.i++
}

func (t *KeyValueTable) Print(w io.Writer) error {

	// Title
	if t.Title != "" {
		fmt.Fprintf(w, "%s\n", color.GreenString(t.Title))
	}

	for i := 0; i < len(t.OrderedData); i++ {
		record := t.OrderedData[i]
		fmt.Fprintf(w, "%-*s: %-*s", t.keyWidth*2, color.BlueString(record.Key), t.valWidth, record.Color(record.Value))
		fmt.Fprintf(w, "\n")
	}
	return nil
}
