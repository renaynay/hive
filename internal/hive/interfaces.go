package hive

import (
	"mime/multipart"
	"net"
)

// Backend captures the docker interactions of the simulation API.
type Backend interface {
	// StartClient starts a client container.
	StartClient(name string, env map[string]string, files map[string]*multipart.FileHeader, checklive bool) (*ClientInfo, error)
	// StopContainer stops the given container.
	StopContainer(containerID string) error

	// RunEnodeSh runs the /enode.sh script in the given container and returns its output.
	RunEnodeSh(containerID string) (string, error)

	// These methods configure docker networks.
	NetworkNameToID(name string) (string, error)
	CreateNetwork(name string) (string, error)
	RemoveNetwork(network string) error
	ContainerIP(containerID, network string) (net.IP, error)
	ConnectContainer(containerID, network string) error
	DisconnectContainer(containerID, network string) error
}
