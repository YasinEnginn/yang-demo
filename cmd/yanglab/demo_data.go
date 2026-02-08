package main

import "yang/internal/models/labnetdevice"

// createDemoData returns the structs for a full network configuration
func createDemoData() (*labnetdevice.Vlans, *labnetdevice.Vrfs, *labnetdevice.Interfaces, *labnetdevice.Routing, *labnetdevice.Bgp, *labnetdevice.System) {
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

	// Interfaces
	enabled := true
	mtu := uint16(1500)

	interfaces := &labnetdevice.Interfaces{
		Interface: []labnetdevice.Interface{
			{
				Name:    "GigabitEthernet0/0",
				Enabled: &enabled,
				Mtu:     &mtu,
				Vrf:     "blue",
				Switchport: &labnetdevice.Switchport{
					Mode:       "access",
					AccessVlan: 10,
				},
				IPv4: &labnetdevice.IPv4{
					Address: []labnetdevice.IPv4Address{
						{IP: "192.0.2.1", PrefixLength: 30},
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
						{IP: "10.0.0.1", PrefixLength: 32},
					},
				},
			},
		},
	}

	// Routing
	routing := &labnetdevice.Routing{
		StaticRoutes: &labnetdevice.StaticRoutes{
			Route: []labnetdevice.StaticRoute{
				{
					Prefix:   "203.0.113.0/24",
					Vrf:      "blue",
					NextHop:  "192.0.2.2",
					Distance: 10,
				},
			},
		},
	}

	// BGP
	bgp := &labnetdevice.Bgp{
		LocalAs: 65001,
		Neighbor: []labnetdevice.Neighbor{
			{
				Address:  "192.0.2.2",
				RemoteAs: 65002,
				Vrf:      "blue",
			},
		},
	}

	return vlans, vrfs, interfaces, routing, bgp, system
}
