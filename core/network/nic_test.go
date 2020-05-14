// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package network_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/core/network"
)

type nicSuite struct {
	info network.InterfaceInfos
}

var _ = gc.Suite(&nicSuite{})

func (s *nicSuite) SetUpTest(_ *gc.C) {
	s.info = network.InterfaceInfos{
		{VLANTag: 1, DeviceIndex: 0, InterfaceName: "eth0"},
		{VLANTag: 0, DeviceIndex: 1, InterfaceName: "eth1"},
		{VLANTag: 42, DeviceIndex: 2, InterfaceName: "br2"},
		{ConfigType: network.ConfigDHCP, NoAutoStart: true},
		{Addresses: network.ProviderAddresses{network.NewProviderAddress("0.1.2.3")}},
		{DNSServers: network.NewProviderAddresses("1.1.1.1", "2.2.2.2")},
		{GatewayAddress: network.NewProviderAddress("4.3.2.1")},
		{AvailabilityZones: []string{"foo", "bar"}},
		{Routes: []network.Route{{
			DestinationCIDR: "0.1.2.3/24",
			GatewayIP:       "0.1.2.1",
			Metric:          0,
		}}},
	}
}

func (s *nicSuite) TestActualInterfaceName(c *gc.C) {
	c.Check(s.info[0].ActualInterfaceName(), gc.Equals, "eth0.1")
	c.Check(s.info[1].ActualInterfaceName(), gc.Equals, "eth1")
	c.Check(s.info[2].ActualInterfaceName(), gc.Equals, "br2.42")
}

func (s *nicSuite) TestIsVirtual(c *gc.C) {
	c.Check(s.info[0].IsVirtual(), jc.IsTrue)
	c.Check(s.info[1].IsVirtual(), jc.IsFalse)
	c.Check(s.info[2].IsVirtual(), jc.IsTrue)
}

func (s *nicSuite) TestIsVLAN(c *gc.C) {
	c.Check(s.info[0].IsVLAN(), jc.IsTrue)
	c.Check(s.info[1].IsVLAN(), jc.IsFalse)
	c.Check(s.info[2].IsVLAN(), jc.IsTrue)
}

func (s *nicSuite) TestAdditionalFields(c *gc.C) {
	c.Check(s.info[3].ConfigType, gc.Equals, network.ConfigDHCP)
	c.Check(s.info[3].NoAutoStart, jc.IsTrue)
	c.Check(s.info[4].Addresses, jc.DeepEquals, network.ProviderAddresses{network.NewProviderAddress("0.1.2.3")})
	c.Check(s.info[5].DNSServers, jc.DeepEquals, network.NewProviderAddresses("1.1.1.1", "2.2.2.2"))
	c.Check(s.info[6].GatewayAddress, jc.DeepEquals, network.NewProviderAddress("4.3.2.1"))
	c.Check(s.info[7].AvailabilityZones, jc.DeepEquals, []string{"foo", "bar"})
	c.Check(s.info[8].Routes, jc.DeepEquals, []network.Route{{
		DestinationCIDR: "0.1.2.3/24",
		GatewayIP:       "0.1.2.1",
		Metric:          0,
	}})
}

func (*nicSuite) TestInterfaceInfosChildren(c *gc.C) {
	interfaces := getInterFaceInfos()

	c.Check(interfaces.Children(""), gc.DeepEquals, interfaces[:2])
	c.Check(interfaces.Children("bond0"), gc.DeepEquals, network.InterfaceInfos{
		interfaces[3], interfaces[4],
	})
	c.Check(interfaces.Children("eth2"), gc.HasLen, 0)
}

func (*nicSuite) TestInterfaceInfosIterHierarchy(c *gc.C) {
	var devs []string
	f := func(info network.InterfaceInfo) error {
		devs = append(devs, info.ParentInterfaceName+":"+info.InterfaceName)
		return nil
	}

	c.Assert(getInterFaceInfos().IterHierarchy(f), jc.ErrorIsNil)

	c.Check(devs, gc.DeepEquals, []string{
		":br-bond0",
		"br-bond0:bond0",
		"bond0:eth0",
		"bond0:eth1",
		":eth2",
	})
}

func getInterFaceInfos() network.InterfaceInfos {
	return network.InterfaceInfos{
		{
			DeviceIndex:   0,
			InterfaceName: "br-bond0",
		},
		{
			DeviceIndex:   1,
			InterfaceName: "eth2",
		},
		{
			DeviceIndex:         2,
			InterfaceName:       "bond0",
			ParentInterfaceName: "br-bond0",
		},
		{
			DeviceIndex:         3,
			InterfaceName:       "eth0",
			ParentInterfaceName: "bond0",
		},
		{
			DeviceIndex:         4,
			InterfaceName:       "eth1",
			ParentInterfaceName: "bond0",
		},
	}
}
