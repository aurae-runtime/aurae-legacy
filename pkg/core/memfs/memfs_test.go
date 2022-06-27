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
	"sync"
	"testing"
)

func TestAddSubNodeSimple(t *testing.T) {
	rootNode.AddSubNode("/test/path", "testData")
	child := rootNode.GetSubNode("/test/path")
	if child.depth != 2 {
		t.Errorf("Expected: 2, Actual: %d", child.depth)
	}
	if child.Name != "path" {
		t.Errorf("Nested child name error. Expected: path, Actual: %s", child.Name)
	}
	if string(child.Content) != "testData" {
		t.Errorf("Nested child value error. Expected: testData, Actual: %s", child.Name)
	}
}

func TestAddSubNodeFileCheck(t *testing.T) {
	rootNode.AddSubNode("/test/path/beeps/boops", "testData")
	child := rootNode.GetSubNode("/test/path/beeps/boops")
	if child == nil {
		t.Errorf("nil child from GetSubNode")
		t.FailNow()
	}
	if child.depth != 4 {
		t.Errorf("Expected: 4, Actual: %d", child.depth)
	}
	if child.Name != "boops" {
		t.Errorf("Nested child name error. Expected: boops, Actual: %s", child.Name)
	}
	if string(child.Content) != "testData" {
		t.Errorf("Nested child value error. Expected: testData, Actual: %s", child.Name)
	}
	if !child.File {
		t.Errorf("Nested child file error, expected child.File=true")
	}
	baseDir := rootNode.GetSubNode("/test/path/beeps")
	if baseDir.File {
		t.Errorf("Base dir file error, expected child.File=false")
	}
	baseDir = rootNode.GetSubNode("/test/path")
	if baseDir.File {
		t.Errorf("Base dir file error, expected child.File=false")
	}
	baseDir = rootNode.GetSubNode("/test")
	if baseDir.File {
		t.Errorf("Base dir file error, expected child.File=false")
	}
	if len(baseDir.Children) != 1 {
		t.Errorf("final basedir check children count failed: %d", len(baseDir.Children))
	}
}

func TestBasicGetSetDepth(t *testing.T) {
	Set("/beeps/boops/meeps/moops", "testvalue")
	result := Get("/beeps/boops/meeps/moops")
	if result != "testvalue" {
		t.Errorf("failed basic test")
	}
}

func TestFuzzCases(t *testing.T) {
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
		Set(c.key, c.expected)
		actual := Get(c.key)
		if actual != c.expected {
			t.Errorf("Expected: %s, Actual: %s", c.expected, actual)
		} else {
			//t.Logf("Expected: %s, Actual: %s", c.expected, actual)
		}
	}

	Set("/beeps/boops/meeps/moops", "testvalue")
	result := Get("/beeps/boops/meeps/moops")
	if result != "testvalue" {
		t.Errorf("failed basic test")
	}
}

func TestListFiles(t *testing.T) {
	Set("/base/path1", "testData1")
	Set("/base/path2", "testData2")
	Set("/base/path3", "testData3")

	actual3 := Get("/base/path3")
	if actual3 != "testData3" {
		t.Errorf("Multiple subfile data lookup error. Expected: testData3, Actual: %s", actual3)
	}
	node := rootNode.GetSubNode("/base/path3")
	if !node.File {
		t.Errorf("Expecting node.File=true")
	}

	files := List("/base")
	if len(files) != 3 {
		t.Errorf("List failure. Expecting 3, Actual: %d", len(files))
	}
	if actual1, ok := files["path1"]; ok {
		if string(actual1.Content) != "testData1" {
			t.Errorf("Expected: testData1, Actual: %s %v", string(actual1.Content), files)
		}
	} else {
		t.Errorf("Unable to find file in list")
	}
	if actual2, ok := files["path2"]; ok {
		if string(actual2.Content) != "testData2" {
			t.Errorf("Expected: testData2, Actual: %s %v", string(actual2.Content), files)
		}
	} else {
		t.Errorf("Unable to find file in list")
	}
	children := rootNode.ListSubNodes("/base")
	for _, node := range children {
		if !node.File {
			t.Errorf("Only expecting files in list")
		}
		if node.depth != 2 {
			t.Errorf("Only expecting 2 depth in list, actual: %d", node.depth)
		}
		if len(node.Children) != 0 {
			t.Errorf("Unexpected sub children.")
		}
	}

}

func TestRemoveNodes(t *testing.T) {
	rootNode.AddSubNode("/ztest/zpath/remove/me", "")
	child := rootNode.GetSubNode("/ztest/zpath")
	childShouldExist := rootNode.GetSubNode("/ztest/zpath/remove/me")
	if childShouldExist == nil {
		t.Errorf("child should exist")
	}
	child.RemoveRecursive()
	childShouldNotExist := rootNode.GetSubNode("/ztest/zpath/remove/me")
	if childShouldNotExist != nil {
		t.Errorf("child should not exist")
	}
	childShouldNoLongerExist := rootNode.GetSubNode("/ztest/zpath")
	if childShouldNoLongerExist != nil {
		t.Errorf("child should no longer exist")
	}
	childShouldStillExist := rootNode.GetSubNode("/ztest")
	if childShouldStillExist == nil {
		t.Errorf("child should still exist")
	}
}

func TestEmptyTree(t *testing.T) {
	rootNode.AddSubNode("/ztest/zpath/remove/me", "")
	rootNode.AddSubNode("/ztest/zzpath/remove/me", "")
	rootNode.AddSubNode("/zzztest/zpath/remove/me", "")
	rootNode.AddSubNode("/zzzzzzzztest/zpath/remove/me", "")
	rootNode.AddSubNode("/zzzzzzztest/zzzzpath/remove/me", "")
	rootNode.RemoveRecursive()
	if rootNode.TotalChildren() != 0 {
		t.Errorf("unable to empty tree")
	}
}

var countChildrenMtx = sync.Mutex{}

func TestCountChildren(t *testing.T) {
	countChildrenMtx.Lock()
	defer countChildrenMtx.Unlock()
	rootNode.RemoveRecursive()
	rootNode.AddSubNode("/single", "")
	if rootNode.TotalChildren() != 1 {
		t.Errorf("Failed tree children count. Expected: 1, Actual: %d", rootNode.TotalChildren())
	}
	rootNode.AddSubNode("/single/double", "")
	if rootNode.TotalChildren() != 2 {
		t.Errorf("Failed tree children count. Expected: 2, Actual: %d", rootNode.TotalChildren())
	}
	rootNode.AddSubNode("/single/double/triple", "")
	if rootNode.TotalChildren() != 3 {
		t.Errorf("Failed tree children count. Expected: 3, Actual: %d", rootNode.TotalChildren())
	}

	doubleNode := rootNode.GetSubNode("/single/double")
	doubleNode.RemoveRecursive()
	if rootNode.TotalChildren() != 1 {
		t.Errorf("Failed tree children count. Expected: 1, Actual: %d", rootNode.TotalChildren())
	}
}
