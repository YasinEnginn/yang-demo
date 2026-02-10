package main

import "yang/internal/models/labnetdevice"

// createDemoData returns the structs for a full network configuration
func createDemoData(profile string) (*labnetdevice.Vlans, *labnetdevice.Vrfs, *labnetdevice.QoS, *labnetdevice.Interfaces, *labnetdevice.Routing, *labnetdevice.Bgp, *labnetdevice.System) {
	isSRLinux := profile == "srlinux"

	// System Users
	system := &labnetdevice.System{
		Users: &labnetdevice.Users{
			User: []labnetdevice.User{
				{UserId: "netadmin", ScreenName: "Network Admin", Role: "admin"},
				{UserId: "operator1", ScreenName: "NOC Operator", Role: "operator"},
			},
		},
	}
	// VLANs
	vlans := &labnetdevice.Vlans{
		Vlan: []labnetdevice.Vlan{
			{Id: 10, Name: "users"},
			{Id: 20, Name: "servers"},
			{Id: 30, Name: "management"},
		},
	}

	// VRFs
	vrfs := &labnetdevice.Vrfs{
		Vrf: []labnetdevice.Vrf{
			{Name: "blue", Rd: "65001:10"},
			{Name: "red", Rd: "65001:20"},
		},
	}

	// QoS Policies
	dscp46 := uint8(46)
	dscp0 := uint8(0)
	bw30 := uint8(30)
	bw50 := uint8(50)

	qos := &labnetdevice.QoS{
		Policy: []labnetdevice.QoSPolicy{
			{
				Name:        "voice-ingress",
				Direction:   "ingress",
				DscpDefault: &dscp46,
				Class: []labnetdevice.QoSClass{
					{
						ClassID:          10,
						ClassName:        "VOICE",
						BandwidthPercent: &bw30,
						PolicingRate:     strPtr("auto"),
					},
				},
			},
			{
				Name:        "wan-egress",
				Direction:   "egress",
				DscpDefault: &dscp0,
				Class: []labnetdevice.QoSClass{
					{
						ClassID:          20,
						ClassName:        "BUSINESS",
						BandwidthPercent: &bw50,
						PolicingRate:     strPtr("100000"), // kbps
					},
				},
			},
		},
	}

	// Interfaces
	enabled := true
	mtu := uint16(1500)

	accessVlan10 := uint16(10)
	pl30 := uint8(30)
	pl32 := uint8(32)

	var accessSwitchport *labnetdevice.Switchport
	if !isSRLinux {
		accessSwitchport = &labnetdevice.Switchport{
			Mode:       "access",
			AccessVlan: &accessVlan10,
		}
	}

	interfaces := &labnetdevice.Interfaces{
		Interface: []labnetdevice.Interface{
			{
				Name:    "GigabitEthernet0/0",
				Enabled: &enabled,
				Mtu:     &mtu,
				Purpose: "lndi:access-port",
				Vrf:     "blue",
				QoS: &labnetdevice.InterfaceQoS{
					InputPolicy:  "voice-ingress",
					OutputPolicy: "wan-egress",
				},
				Switchport: accessSwitchport,
				IPv4: &labnetdevice.IPv4{
					Address: []labnetdevice.IPv4Address{
						{IP: "192.0.2.1", PrefixLength: &pl30},
					},
				},
			},
			{
				Name:    "Loopback0",
				Enabled: &enabled,
				Mtu:     &mtu,
				Vrf:     "red",
				IPv4: &labnetdevice.IPv4{
					Address: []labnetdevice.IPv4Address{
						{IP: "10.0.0.1", PrefixLength: &pl32},
					},
				},
			},
		},
	}

	// Routing
	nh := "192.0.2.2"
	dist10 := uint8(10)

	routing := &labnetdevice.Routing{
		StaticRoutes: &labnetdevice.StaticRoutes{
			Route: []labnetdevice.StaticRoute{
				{
					Prefix:   "203.0.113.0/24",
					Vrf:      "blue",
					NextHop:  &nh,
					Distance: &dist10,
				},
			},
		},
	}

	// BGP
	localAs := uint32(65001)
	remoteAs := uint32(65002)
	bgpVrf := "blue"
	if isSRLinux {
		bgpVrf = ""
	}

	bgp := &labnetdevice.Bgp{
		LocalAs: &localAs,
		Neighbor: []labnetdevice.Neighbor{
			{
				Address:  "192.0.2.2",
				RemoteAs: &remoteAs,
				Vrf:      bgpVrf,
			},
		},
	}

	return vlans, vrfs, qos, interfaces, routing, bgp, system
}

func strPtr(s string) *string {
	return &s
}
