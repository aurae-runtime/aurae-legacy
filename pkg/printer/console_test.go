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
	c.AddPrinter(t1)
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

	c.AddPrinter(t1)
	err := c.PrintStdout()
	if err != nil {
		t.Errorf("unable to print: %v", err)
	}
}

func TestKeyValue(t *testing.T) {

	c := NewConsole("TestConsoleDoubleField")
	t1 := NewKeyValueTable("") // No title does not print

	t1.AddKeyValue("Beeps", "Boops")
	t1.AddKeyValue("Meeps", 1)
	t1.AddKeyValue("Boops", &struct {
		something string
	}{
		something: "anything",
	})

	c.AddPrinter(t1)
	err := c.PrintStdout()
	if err != nil {
		t.Errorf("unable to print: %v", err)
	}
}

func TestMultiPrint(t *testing.T) {

	c := NewConsole("TestMultiPrint")
	t1 := NewKeyValueTable("") // No title does not print

	t1.AddKeyValue("Beeps", "Boops")
	t1.AddKeyValue("Meeps", 1)
	t1.AddKeyValue("Boops", &struct {
		something string
	}{
		something: "anything",
	})

	c.AddPrinter(t1)
	t2 := NewTable("")

	field := t2.NewField("boops")
	field.AddValue("meeps")
	field.AddValue("moops")
	field.AddValue("sheesh")
	t2.AddField(field)

	field = t2.NewField("beeps")
	field.AddValue("meeps")
	field.AddValue("moops")
	field.AddValue("zeeps")
	t2.AddField(field)

	field = t2.NewField("meeps")
	field.AddValue("zerks")
	field.AddValue("jeepers")
	field.AddValue("zeeps")
	field.AddValue("zeeks")
	t2.AddField(field)

	field = t2.NewField("jeepers")
	field.AddValue("zorps")
	field.AddValue("borks")
	field.AddValue("zeeps")
	field.AddValue("zeeks")
	t2.AddField(field)

	field = t2.NewField("jenkies")
	field.AddValue("zorps")
	field.AddValue("zoops")
	field.AddValue("eeks")
	field.AddValue("zeeks")
	t2.AddField(field)

	field = t2.NewField("zeeps")
	field.AddValue("zorps")
	field.AddValue("zoops")
	field.AddValue("zeeps")
	field.AddValue("zeeks")
	t2.AddField(field)

	c.AddPrinter(t1)
	c.AddPrinter(t2)

	err := c.PrintStdout()
	if err != nil {
		t.Errorf("unable to print: %v", err)
	}
}
