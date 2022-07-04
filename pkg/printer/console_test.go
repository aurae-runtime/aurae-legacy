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

import "testing"

func TestConsoleSingleField(t *testing.T) {

	c := NewConsole("TestConsoleSingleField")
	t1 := NewTable("People")
	nameField := t1.NewField("Names")
	nameField.AddValue("Kris")
	nameField.AddValue("Nóva")
	nameField.AddValue("Björn")
	t1.AddField(nameField)
	c.AddTable(t1)
	err := c.PrintStdout()
	if err != nil {
		t.Errorf("unable to print: %v", err)
	}
}

func TestConsoleDoubleField(t *testing.T) {

	c := NewConsole("TestConsoleDoubleField")
	t1 := NewTable("People")

	nameField := t1.NewField("Names")
	nameField.AddValue("Kris")
	nameField.AddValue("Quintessence")
	nameField.AddValue("Björn")
	t1.AddField(nameField)

	thingField := t1.NewField("Favorite Things")
	thingField.AddValue("Mountains")
	thingField.AddValue("Stars")
	thingField.AddValue("Gravy")
	t1.AddField(thingField)

	c.AddTable(t1)
	err := c.PrintStdout()
	if err != nil {
		t.Errorf("unable to print: %v", err)
	}
}
