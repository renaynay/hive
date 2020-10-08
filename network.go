package main

import (
	"errors"
	"fmt"
	"net"
	"strings"

	docker "github.com/fsouza/go-dockerclient"

	"gopkg.in/inconshreveable/log15.v2"
)

func createNetworks(networkNames []string) ([]*docker.Network, error) {
	// list networks to make sure not to duplicate
	existing, err := dockerClient.ListNetworks()
	if err != nil {
		return nil, err
	}
	// range over existing networks and if one of their names matches a specified network name,
	// remove it
	for _, exists := range existing {
		for _, name := range networkNames {
			if exists.Name == name {
				if err := dockerClient.RemoveNetwork(exists.ID); err != nil {
					return nil, err
				}
			}
		}
	}
	// create networks
	networks := make([]*docker.Network, 0)
	for _, name := range networkNames {
		network, err := dockerClient.CreateNetwork(docker.CreateNetworkOptions{
			Name:           name,
			CheckDuplicate: true, // TODO set this to tru?
			Attachable:     true,
		})
		if err != nil {
			return nil, err
		}

		networks = append(networks, network)
	}
	// if no networks were created, abort
	if len(networks) == 0 {
		return nil, fmt.Errorf("no networks created")
	}
	return networks, nil
}

func connectContainer(network, container string) error {
	return dockerClient.ConnectNetwork(network, docker.NetworkConnectionOptions{
		Container:      container,
		EndpointConfig: nil, // TODO ?
	})
}

// lookupBridgeIP attempts to locate the IPv4 address of the local docker0 bridge
// network adapter.
func lookupBridgeIP(logger log15.Logger) (net.IP, error) {
	// Find the local IPv4 address of the docker0 bridge adapter
	interfaes, err := net.Interfaces()
	if err != nil {
		logger.Error("failed to list network interfaces", "error", err)
		return nil, err
	}
	// Iterate over all the interfaces and find the docker0 bridge
	for _, iface := range interfaes {
		if iface.Name == "docker0" || strings.Contains(iface.Name, "vEthernet") {
			// Retrieve all the addresses assigned to the bridge adapter
			addrs, err := iface.Addrs()
			if err != nil {
				logger.Error("failed to list docker bridge addresses", "error", err)
				return nil, err
			}
			// Find a suitable IPv4 address and return it
			for _, addr := range addrs {
				ip, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					logger.Error("failed to list parse address", "address", addr, "error", err)
					return nil, err
				}
				if ipv4 := ip.To4(); ipv4 != nil {
					return ipv4, nil
				}
			}
		}
	}
	// Crap, no IPv4 found, bounce
	return nil, errors.New("not found")
}
