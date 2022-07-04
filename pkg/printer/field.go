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
	"strings"
)

type Field struct {
	Header string
	i      int
	values map[int]string
	width  int
}

func (t *Table) NewField(header string) *Field {
	header = fmt.Sprintf("%s%s%s", FieldPaddingLeft, strings.ToUpper(header), FieldPaddingRight)
	return &Field{
		Header: header,
		i:      0,
		values: make(map[int]string),
		width:  0,
	}
}

const (
	FieldPaddingLeft  string = " "
	FieldPaddingRight string = " "
)

func (f *Field) AddValue(v any) {
	str := fmt.Sprintf("%s%v%s", FieldPaddingLeft, v, FieldPaddingRight)
	if len(str) > f.width {
		f.width = len(str)
	}
	f.values[f.i] = str
	f.i++
}
