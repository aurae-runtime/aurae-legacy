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

package core

import (
	"context"
	"github.com/kris-nova/aurae/rpc"
	"strings"
	"testing"
)

func TestBasicIOSad(t *testing.T) {

	db := NewPathDatabase()

	// Set
	var setResp *rpc.SetResp
	setResp, err := db.SetRPC(context.Background(), &rpc.SetReq{
		Key: "",
		Val: "testBadData",
	})
	if err != nil {
		t.Errorf("unable to SetRPC: %v", err)
	}
	if setResp.Code != CoreCode_EMPTY {
		t.Errorf("Invalid response code. Expected: %d, Actual: %d", CoreCode_OKAY, setResp.Code)
	}
}

func TestBasicIOHappy(t *testing.T) {

	db := NewPathDatabase()

	// Set
	var setResp *rpc.SetResp
	setResp, err := db.SetRPC(context.Background(), &rpc.SetReq{
		Key: "testKey",
		Val: "testVal",
	})
	if err != nil {
		t.Errorf("unable to SetRPC: %v", err)
	}
	if setResp.Code != CoreCode_OKAY {
		t.Errorf("Invalid response code. Expected: %d, Actual: %d", CoreCode_OKAY, setResp.Code)
	}

	// Get
	var getResp *rpc.GetResp
	getResp, err = db.GetRPC(context.Background(), &rpc.GetReq{
		Key: "testKey",
	})
	if err != nil {
		t.Errorf("unable to GetRPC: %v", err)
	}
	if getResp.Val != "testVal" {
		t.Errorf("Database IO inconsistency. Expected: %s, Actual: %s", "testVal", getResp.Val)
	}
	if getResp.Code != CoreCode_OKAY {
		t.Errorf("Invalid response code. Expected: %d, Actual: %d", CoreCode_OKAY, getResp.Code)
	}

}

func TestIOCases(t *testing.T) {

	db := NewPathDatabase()

	cases := []struct {
		getKey       string
		setKey       string
		setVal       string
		expected     string
		expectedCode int32
	}{
		{
			// Basic test
			setKey:       "beeps",
			setVal:       "boops",
			getKey:       "beeps",
			expected:     "boops",
			expectedCode: CoreCode_OKAY,
		},
		{
			// Ensure the path was created by memfs
			setKey:       "beeps",
			setVal:       "boops",
			getKey:       "/beeps",
			expected:     "boops",
			expectedCode: CoreCode_OKAY,
		},
		{
			// Ensure the path was formatted by memfs
			setKey:       "//////beeps",
			setVal:       "boops",
			getKey:       "/beeps",
			expected:     "boops",
			expectedCode: CoreCode_OKAY,
		},
		{
			// Ensure the path was created by memfs complex
			setKey:       "/////\\/beeps/example",
			setVal:       "boops",
			getKey:       "/beeps/example",
			expected:     "boops",
			expectedCode: CoreCode_OKAY,
		},
		{
			// Ensure the path was created by memfs
			setKey:       "**\\/\\beeps",
			setVal:       "boops",
			getKey:       "/**/beeps",
			expected:     "boops",
			expectedCode: CoreCode_OKAY,
		},
	}

	for _, c := range cases {
		// Set
		var setResp *rpc.SetResp
		setResp, err := db.SetRPC(context.Background(), &rpc.SetReq{
			Key: c.setKey,
			Val: c.setVal,
		})
		if setResp.Code != CoreCode_OKAY {
			t.Errorf("Assumed OKAY. Actual: %d", setResp.Code)
		}
		if err != nil {
			t.Errorf("unable to SetRPC: %v", err)
		}
		// Get
		var getResp *rpc.GetResp
		getResp, err = db.GetRPC(context.Background(), &rpc.GetReq{
			Key: c.getKey,
		})
		if err != nil {
			t.Errorf("unable to GetRPC: %v", err)
		}
		if getResp.Val != c.expected {
			t.Errorf("Unexpected data IO: Expected: %s, Actual: %s", c.expected, getResp.Val)
		}
		if getResp.Code != c.expectedCode {
			t.Errorf("Unexpected code: Expected: %d, Actual: %d", c.expectedCode, getResp.Code)
		}
	}
}

func TestBasicListIOHappy(t *testing.T) {

	db := NewPathDatabase()

	// Set
	var setResp *rpc.SetResp
	setResp, err := db.SetRPC(context.Background(), &rpc.SetReq{
		Key: "testKey",
		Val: "testVal",
	})
	if err != nil {
		t.Errorf("unable to SetRPC: %v", err)
	}
	if setResp.Code != CoreCode_OKAY {
		t.Errorf("Invalid response code. Expected: %d, Actual: %d", CoreCode_OKAY, setResp.Code)
	}

	// List
	var lsResp *rpc.ListResp
	lsResp, err = db.ListRPC(context.Background(), &rpc.ListReq{
		Key: "testKey",
	})
	if err != nil {
		t.Errorf("unable to GetRPC: %v", err)
	}
	if lsResp.Code != CoreCode_OKAY {
		t.Errorf("Invalid response code. Expected: %d, Actual: %d", CoreCode_OKAY, lsResp.Code)
	}

	if dirent, ok := lsResp.Entries["testKey"]; !ok {
		t.Errorf("Missing testKey in list")
	} else {
		if dirent.Name != "testKey" {
			t.Errorf("List data IO. Expected: %s, Actual: %s", "testKey", dirent.Name)
		}
		if dirent.File != true {
			t.Errorf("Expected file=true, Actual file=false")
		}
	}
}

func TestComplexListIOCases(t *testing.T) {

	db := NewPathDatabase()

	cases := []struct {
		setKeys          []string
		setVal           string
		listKey          string
		expectedListKeys []string
	}{
		{
			// Set a file called "/test/path" to value "testVal" and ensure
			// list on "/test" returns the file "path"
			setKeys:          []string{"test/path"},
			setVal:           "testVal",
			listKey:          "test",
			expectedListKeys: []string{"path"},
		},
		{
			// Check 2 files, one with a leading slash, one without
			setKeys:          []string{"test/path1", "/test/path2"},
			setVal:           "testVal",
			listKey:          "test",
			expectedListKeys: []string{"path1", "path2"},
		},
		{
			// Ensure only /test returns the expected files by adding bad files
			setKeys:          []string{"test/path1", "/test/path2", "beeps/boops", "bad/path"},
			setVal:           "testVal",
			listKey:          "test",
			expectedListKeys: []string{"path1", "path2"},
		},
		{
			// Basic nested file test
			setKeys:          []string{"test/path/file1", "/test/path/file2"},
			setVal:           "testVal",
			listKey:          "/test/path",
			expectedListKeys: []string{"file1", "file2"},
		},
		{
			// Check if a node is added that should already exist we still return
			// the correct results
			setKeys:          []string{"test/path/file1", "/test/path/file2", "test", "test/path", "/test/path", "/test"},
			setVal:           "testVal",
			listKey:          "/test/path",
			expectedListKeys: []string{"file1", "file2"},
		},
		{
			// Ensure that files can be changed to dirs as keys are nested under them
			setKeys:          []string{"/dir1/dir2/file1", "dir1/dir2/file2", "dir1/dir2/file2/fileX"},
			setVal:           "testVal",
			listKey:          "/dir1/dir2/file2",
			expectedListKeys: []string{"fileX"},
		},
		{
			// Sad test to ensure we aren't polluting the data
			setKeys:          []string{"/dir1/dir2/file1", "dir1/dir2/file2", "dir1/dir2/file2/fileX"},
			setVal:           "testVal",
			listKey:          "/bad/path",
			expectedListKeys: []string{""},
		},
		{
			// Ensure we list sub directories
			setKeys:          []string{"/dir1/dir2/file1", "dir1/dir2/file2", "dir1/dir2/file2/fileX"},
			setVal:           "testVal",
			listKey:          "/dir1",
			expectedListKeys: []string{"dir2"},
		},
		{
			// Ensure we list multiple sub directories
			setKeys:          []string{"/dir1/dirSub1/file1", "dir1/dirSub2/file1"},
			setVal:           "testVal",
			listKey:          "/dir1",
			expectedListKeys: []string{"dirSub1", "dirSub2"},
		},
	}

	for _, c := range cases {
		// Set
		for _, setKey := range c.setKeys {
			var setResp *rpc.SetResp
			setResp, err := db.SetRPC(context.Background(), &rpc.SetReq{
				Key: setKey,
				Val: c.setVal,
			})
			if err != nil {
				t.Errorf("unable to SetRPC: %v", err)
			}
			if setResp.Code != CoreCode_OKAY {
				t.Errorf("Invalid response code. Expected: %d, Actual: %d", CoreCode_OKAY, setResp.Code)
			}
		}

		// List
		var lsResp *rpc.ListResp
		lsResp, err := db.ListRPC(context.Background(), &rpc.ListReq{
			Key: c.listKey,
		})
		if err != nil {
			t.Errorf("unable to GetRPC: %v", err)
		}
		if lsResp.Code != CoreCode_OKAY {
			t.Errorf("Invalid response code. Expected: %d, Actual: %d", CoreCode_OKAY, lsResp.Code)
		}
		for _, expectedKey := range c.expectedListKeys {
			if expectedKey == "" {
				continue
			}
			if dirent, ok := lsResp.Entries[expectedKey]; !ok {
				t.Errorf("Missing %s in list", expectedKey)
				t.Errorf("Returned keys: %s", strings.Join(listResponseToStrings(lsResp), " "))
			} else {
				if dirent.Name != expectedKey {
					t.Errorf("List data IO. Expected: %s, Actual: %s", "testKey", dirent.Name)
				}
			}
		}

	}

}

func listResponseToStrings(lsResp *rpc.ListResp) []string {
	var ret []string
	for k, _ := range lsResp.Entries {
		ret = append(ret, k)
	}
	return ret
}

//func TestTODO(t *testing.T) {
//
//	db := NewPathDatabase()
//
//	// Set
//	var setResp *rpc.SetResp
//	setResp, err := db.SetRPC(context.Background(), &rpc.SetReq{})
//	if err != nil {
//
//	}
//
//	// Get
//	var getResp *rpc.GetResp
//	getResp, err = db.GetRPC(context.Background(), &rpc.GetReq{})
//	if err != nil {
//
//	}
//
//}
