package endpoints

import "fmt"

type SocketEndpoint string
type ExternalSocketEndpoint string

//Purpose of this package is to provide an overview of all socket endpoints and their binding scheme

const (
	// COLLECTOR
	//PUSH/PULL binds on collector worker
	COLLECTOR SocketEndpoint = "collector"

	// STORAGE
	//PUSH/PULL binds on storage worker
	STORAGE SocketEndpoint = "storage"

	// STORAGE_PROVIDE
	//REQ/REP binds on storage worker
	STORAGE_PROVIDE SocketEndpoint = "storage_provide"

	// DETECTOR
	//PUSH/PULL binds on detector worker
	DETECTOR SocketEndpoint = "detection"
)

var (
	// STORAGE_API
	//REQ/REP binds on storage worker
	//Like
	STORAGE_API ExternalSocketEndpoint
)

func init() {
	STORAGE_API = ExternalSocketEndpoint(
		fmt.Sprintf("127.0.0.1:%d", 5555)) //TODO: read from environment
}

// InProcessEndpoint returns an endpoint string using the inproc:// protocol
func InProcessEndpoint(endpoint SocketEndpoint) string {
	return fmt.Sprintf("inproc://%s", endpoint)
}

// TcpEndpoint returns an endpoint string using the tcp:// protocol
func TcpEndpoint(endpoint ExternalSocketEndpoint) string {
	return fmt.Sprintf("tcp://%s", endpoint)
}
