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

package aurafs

import (
	"encoding/json"
	"fmt"
	"sync"
)

const (
	FileauraeStatus string = "status"
)

type StatusValues struct {
	Value     string
	Socket    string
	Listening bool
}

var (
	internalStatusMtx = sync.Mutex{}
	internalStatus    = &StatusValues{}
)

func GetStatus() *StatusValues {
	internalStatusMtx.Lock()
	defer internalStatusMtx.Unlock()
	return internalStatus
}

func SetStatus(s *StatusValues) {
	internalStatusMtx.Lock()
	defer internalStatusMtx.Unlock()
	internalStatus = s
}

// GetStatusJsonBytes will render the status to JSON and return the
// in-memory data with a newline at the end of the file.
func GetStatusJsonBytes() []byte {
	internalStatusMtx.Lock()
	defer internalStatusMtx.Unlock()
	var data []byte
	data, _ = json.Marshal(internalStatus)
	fileData := fmt.Sprintf("%s\n", string(data))
	return []byte(fileData)
}
