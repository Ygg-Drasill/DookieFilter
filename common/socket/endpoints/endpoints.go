package endpoints

import "fmt"

type SocketEndpoint string
type ExternalSocketEndpoint string

//Purpose of this package is to provide an overview of all socket endpoints and their binding scheme

const (
	// STORAGE
	//PUSH/PULL binds on storage worker
	STORAGE SocketEndpoint = "storage"

	// STORAGE_PROVIDE
	//REQ/REP binds on storage worker
	STORAGE_PROVIDE SocketEndpoint = "storage_provide"

	// DETECTOR
	//PUSH/PULL binds on detector worker
	DETECTOR SocketEndpoint = "detection"

	// FILTER
	//PUSH/PULL binds on filter worker
	FILTER_INPUT SocketEndpoint = "filter"
)

var (
	// COLLECTOR
	//PUSH/PULL binds on collector worker
	COLLECTOR ExternalSocketEndpoint
	// STORAGE_API
	//REQ/REP binds on storage worker
	//Like
	STORAGE_API   ExternalSocketEndpoint
	FILTER_OUTPUT ExternalSocketEndpoint
	IMPUTATION    ExternalSocketEndpoint
	STORAGE_PROXY ExternalSocketEndpoint
)

func init() {
	COLLECTOR = ExternalSocketEndpoint(
		fmt.Sprintf("127.0.0.1:%d", 5559))
	STORAGE_API = ExternalSocketEndpoint(
		fmt.Sprintf("127.0.0.1:%d", 5560)) //TODO: read from environment
	FILTER_OUTPUT = ExternalSocketEndpoint(
		fmt.Sprintf("127.0.0.1:%d", 5556)) //TODO: read from environment
	IMPUTATION = ExternalSocketEndpoint(
		fmt.Sprintf("127.0.0.1:%d", 5557)) //TODO: read from environment
	STORAGE_PROXY = ExternalSocketEndpoint(
		fmt.Sprintf("127.0.0.1:%d", 5558)) //TODO: read from environment
}

// InProcessEndpoint returns an endpoint string using the inproc:// protocol
func InProcessEndpoint(endpoint SocketEndpoint) string {
	return fmt.Sprintf("inproc://%s", endpoint)
}

// TcpEndpoint returns an endpoint string using the tcp:// protocol
func TcpEndpoint(endpoint ExternalSocketEndpoint) string {
	return fmt.Sprintf("tcp://%s", endpoint)
}
