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

package main

import (
	"context"
	"fmt"
	"github.com/kris-nova/aurae/client"
	"github.com/kris-nova/aurae/gen/aurae"
	"github.com/kris-nova/aurae/pkg/common"
	"github.com/kris-nova/aurae/pkg/daemon"
	"os"
	"strings"
)

func runtime(args []string) error {

	if len(args) < 1 {
		return fmt.Errorf("missing arguments to execute")
	}

	ctx := context.Background()

	x := client.NewClient()
	err := x.ConnectSocket(daemon.DefaultSocketLocationLinux)
	if err != nil {
		return err
	}

	// Tell Aurae to leverage the exec service
	x.AdoptService(ctx, &aurae.AdoptServiceRequest{
		UniqueComponentName: "exec", // Adopt the exec service in /providers/service/exec
	})
	// TODO Error type "Already registered"
	//if err != nil {
	//	return fmt.Errorf("unable to adopt service: %v", err)
	//}
	//if adoptResp.Code != common.ResponseCode_OKAY {
	//	return fmt.Errorf("unable to adopt service: %v", adoptResp.Message)
	//}

	// Exec the arguments
	runResp, err := x.RunProcess(ctx, &aurae.RunProcessRequest{
		ExecutableCommand: strings.Join(args, " "),
	})
	if err != nil {
		return fmt.Errorf("unable to run process: %v", err)
	}
	if runResp.Code != common.ResponseCode_OKAY {
		return fmt.Errorf("unable to run process: %v", runResp.Message)
	}
	pid := runResp.PID

	readResp, err := x.ReadStdout(ctx, &aurae.ReadStdoutRequest{
		PID:    pid,
		Length: 1024,
	})
	if err != nil {
		return fmt.Errorf("unable to read stdout: %v", err)
	}
	if readResp.Code != common.ResponseCode_OKAY {
		return fmt.Errorf("unable to read stdout: %v", readResp.Message)
	}

	fmt.Println(readResp.Data)

	return nil
}

func main() {
	err := runtime(os.Args[1:])
	if err != nil {
		fmt.Println("Error running example:")
		fmt.Printf("%+v\n", err)
	}
}
