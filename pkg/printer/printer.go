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
	"reflect"
)

func PrintStdout(title string, v any) {
	c := AnyToConsole(title, v)
	c.PrintStdout()
}

func PrintStderr(title string, v any) {
	c := AnyToConsole(title, v)
	c.PrintStderr()
}

func AnyToConsole(title string, v interface{}) *Console {
	c := NewConsole(title)
	r := reflect.ValueOf(v).Elem()
	if r.Kind() == reflect.Ptr {
		r = reflect.Indirect(r)
	}
	rType := r.Type()

	// HasLists
	hasLists := false
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		if f.Kind() == reflect.Map {
			hasLists = true
			break
		}
		if f.Kind() == reflect.Slice {
			hasLists = true
		}
	}

	if hasLists {
		// Table Mode
		t := NewTable("")
		for i := 0; i < r.NumField(); i++ {
			f := r.Field(i)
			pf := t.NewField(rType.Field(i).Name)
			if f.Kind() == reflect.Map {
				if len(f.MapKeys()) == 0 {
					pf.AddValue("---------")
				} else {
					for _, k := range f.MapKeys() {
						v := f.MapIndex(k)
						pf.AddValue(fmt.Sprintf("%s:%s", k, v))
					}
				}
			} else {
				pf.AddValue(f.Interface())
			}
			t.AddField(pf)

		}
		c.AddPrinter(t)
	} else {
		// Key Value Mode
		kv := NewKeyValueTable("")
		for i := 0; i < r.NumField(); i++ {
			f := r.Field(i)
			if f.CanInterface() {
				kv.AddKeyValue(fmt.Sprintf("%v %v", rType.Field(i).Name, f.Type()), f.Interface())
			}
		}
		c.AddPrinter(kv)
	}
	return c
}
