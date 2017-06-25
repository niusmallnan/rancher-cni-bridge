package metadata

import (
	"net"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher-metadata/metadata"
)

const (
	metadataURL         = "http://rancher-metadata/2015-12-19"
	multiplierForTwoMin = 240
	subnetLabel         = "io.rancher.vpc.subnet"
	emptySubnet         = ""
)

type VPCFinderFromMetadata struct {
	m metadata.Client
}

func NewVPCFinderFromMetadata() (*VPCFinderFromMetadata, error) {
	m := metadata.NewClient(metadataURL)
	return &VPCFinderFromMetadata{m}, nil
}

func (vf *VPCFinderFromMetadata) GetSelfSubnet() string {
	for i := 0; i < multiplierForTwoMin; i++ {
		host, err := vf.m.GetSelfHost()
		if err != nil {
			logrus.Errorf("rancher-cni-bridge: Error getting metadata host: %v", err)
			return emptySubnet
		}

		cidr, ok := host.Labels[subnetLabel]
		if ok {
			_, _, err = net.ParseCIDR(cidr)
			if err != nil {
				logrus.Errorf("rancher-cni-bridge: Invalid vpc subnet: %s", cidr)
				return emptySubnet
			}
			return cidr
		}

		logrus.Info("rancher-cni-bridge: Waiting to find the vpc subnet host key")
		time.Sleep(500 * time.Millisecond)
	}
	logrus.Info("rancher-cni-bridge: VPC subnet host key not found: %s ", subnetLabel)
	return emptySubnet
}
