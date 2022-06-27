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
		t.Errorf("Invalid response code. Expected: %d, Actual: %d", CoreCode_OKAY, getResp.Code)
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
