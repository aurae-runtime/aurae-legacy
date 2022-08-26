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

package config

import (
	"context"
	"github.com/kris-nova/aurae/gen/aurae"
"github.com/kris-nova/aurae/pkg/config/local"
"testing"
)

const localStateBase = "/tmp/aurae.test"

func TestBasicIOLocalSad(t *testing.T) {
	getFromMemory = false
	stateStore := local.NewState(localStateBase)
	db := NewService(stateStore)

	// Set
	var setResp *rpc.SetResp
	setResp, err := db.Set(context.Background(), &rpc.SetReq{
		Key: "",
		Val: "testBadData",
	})
	if err != nil {
		t.Errorf("unable to Set: %v", err)
	}
	if setResp.Code != CoreCode_EMPTY {
		t.Errorf("Invalid response code. Expected: %d, Actual: %d", CoreCode_OKAY, setResp.Code)
	}
}

func TestBasicIOLocalHappy(t *testing.T) {
	getFromMemory = false

	stateStore := local.NewState(localStateBase)
	db := NewService(stateStore)

	// Set
	var setResp *rpc.SetResp
	setResp, err := db.Set(context.Background(), &rpc.SetReq{
		Key: "testKey",
		Val: "testVal",
	})
	if err != nil {
		t.Errorf("unable to Set: %v", err)
	}
	if setResp.Code != CoreCode_OKAY {
		t.Errorf("Invalid response code. Expected: %d, Actual: %d", CoreCode_OKAY, setResp.Code)
	}

	// Get
	var getResp *rpc.GetResp
	getResp, err = db.Get(context.Background(), &rpc.GetReq{
		Key: "testKey",
	})
	if err != nil {
		t.Errorf("unable to Get: %v", err)
	}
	if getResp.Val != "testVal" {
		t.Errorf("Database IO inconsistency. Expected: %s, Actual: %s", "testVal", getResp.Val)
	}
	if getResp.Code != CoreCode_OKAY {
		t.Errorf("Invalid response code. Expected: %d, Actual: %d", CoreCode_OKAY, getResp.Code)
	}

}
