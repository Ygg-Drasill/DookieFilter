package endpoints

import "fmt"

type SocketEndpoint string

const (
	// COLLECTOR PUSH/PULL binds on collector worker
	COLLECTOR SocketEndpoint = "collector"
	// STORAGE PUSH/PULL binds on storage worker
	STORAGE SocketEndpoint = "storage"
	// STORAGE_PROVIDE REQ/REP binds on storage worker
	STORAGE_PROVIDE SocketEndpoint = "storage_provide"
	// DETECTOR PUSH/PULL binds on detector worker
	DETECTOR SocketEndpoint = "detection"
)

func InProcessEndpoint(endpoint SocketEndpoint) string {
	return fmt.Sprintf("inproc://collector/%s", endpoint)
}
