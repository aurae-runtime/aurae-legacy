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

package memfs

import (
	"testing"
)

func TestAddChildSimple(t *testing.T) {
	rootNode.AddChild("/test/path", "testData")
	child := rootNode.GetChild("/test/path")
	if child.depth != 3 {
		t.Errorf("Expected: 3, Actual: %d", child.depth)
	}
	if child.Name != "path" {
		t.Errorf("Nested child name error. Expected: path, Actual: %s", child.Name)
	}
	if child.Value != "testData" {
		t.Errorf("Nested child value error. Expected: testData, Actual: %s", child.Name)
	}
}

func TestAddChildFileCheck(t *testing.T) {
	rootNode.AddChild("/test/path/beeps/boops", "testData")
	child := rootNode.GetChild("/test/path/beeps/boops")
	if child == nil {
		t.Errorf("nil child from GetChild")
		t.FailNow()
	}
	if child.depth != 5 {
		t.Errorf("Expected: 5, Actual: %d", child.depth)
	}
	if child.Name != "boops" {
		t.Errorf("Nested child name error. Expected: boops, Actual: %s", child.Name)
	}
	if child.Value != "testData" {
		t.Errorf("Nested child value error. Expected: testData, Actual: %s", child.Name)
	}
	if !child.file {
		t.Errorf("Nested child file error, expected child.file=true")
	}
	baseDir := rootNode.GetChild("/test/path/beeps")
	if baseDir.file {
		t.Errorf("Base dir file error, expected child.file=false")
	}
	baseDir = rootNode.GetChild("/test/path")
	if baseDir.file {
		t.Errorf("Base dir file error, expected child.file=false")
	}
	baseDir = rootNode.GetChild("/test")
	if baseDir.file {
		t.Errorf("Base dir file error, expected child.file=false")
	}
	if len(baseDir.Children) != 1 {
		t.Errorf("final basedir check children count failed: %d", len(baseDir.Children))
	}
}

func TestBasicGetSetDepth(t *testing.T) {
	db := NewDatabase()
	db.Set("/beeps/boops/meeps/moops", "testvalue")
	result := db.Get("/beeps/boops/meeps/moops")
	if result != "testvalue" {
		t.Errorf("failed basic test")
	}
}

func TestFuzzCases(t *testing.T) {
	db := NewDatabase()
	cases := []struct {
		key      string
		expected string
	}{
		{
			key:      "boops",
			expected: "/boops",
		},
		{
			key:      "boops///",
			expected: "/boops",
		},
		{
			key:      "//boops",
			expected: "/boops",
		},
		{
			key:      "//\\/\\/\\//\\/\\//boops",
			expected: "/boops",
		},
		{
			key:      "beeps/boops/  zeeps",
			expected: "/beeps/boops/zeeps",
		},
	}

	for _, c := range cases {
		db.Set(c.key, c.expected)
		actual := db.Get(c.key)
		if actual != c.expected {
			t.Errorf("Expected: %s, Actual: %s", c.expected, actual)
		} else {
			//t.Logf("Expected: %s, Actual: %s", c.expected, actual)
		}
	}

	db.Set("/beeps/boops/meeps/moops", "testvalue")
	result := db.Get("/beeps/boops/meeps/moops")
	if result != "testvalue" {
		t.Errorf("failed basic test")
	}
}

func TestListFiles(t *testing.T) {
	db := NewDatabase()
	db.Set("/base/path1", "testData1")
	db.Set("/base/path2", "testData2")
	db.Set("/base/path3", "testData3")

	actual3 := db.Get("/base/path3")
	if actual3 != "testData3" {
		t.Errorf("Multiple subfile data lookup error. Expected: testData3, Actual: %s", actual3)
	}
	node := rootNode.GetChild("/base/path3")
	if !node.file {
		t.Errorf("Expecting node.file=true")
	}

	files := db.List("/base")
	if len(files) != 3 {
		t.Errorf("List failure. Expecting 3, Actual: %d", len(files))
	}
	if actual1, ok := files["path1"]; ok {
		if actual1 != "testData1" {
			t.Errorf("Expected: testData1, Actual: %s %v", actual1, files)
		}
	} else {
		t.Errorf("Unable to find file in list")
	}
	if actual2, ok := files["path2"]; ok {
		if actual2 != "testData2" {
			t.Errorf("Expected: testData2, Actual: %s %v", actual2, files)
		}
	} else {
		t.Errorf("Unable to find file in list")
	}
	children := rootNode.ListChildren("/base")
	for _, node := range children {
		if !node.file {
			t.Errorf("Only expecting files in list")
		}
		if node.depth != 3 {
			t.Errorf("Only expecting 3 depth in list, actual: %d", node.depth)
		}
		if len(node.Children) != 0 {
			t.Errorf("Unexpected sub children.")
		}
	}

}
