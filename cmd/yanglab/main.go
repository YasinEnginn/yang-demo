package main

import (
	"fmt"
	"log"
	"yang/internal/client"
	"yang/internal/models/labnetdevice"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("       YANG LAB - NETWORK MANAGER      ")
	fmt.Println("========================================")

	// 1. Connect
	c, err := client.New("127.0.0.1:830", "netconf", "netconf")
	if err != nil {
		log.Fatalf("[-] Connection Failed: %v", err)
	}
	defer c.Close()
	fmt.Println("[+] Connected to NETCONF Server (127.0.0.1:830)")

	// 2. Execute Operations directly
	pushNetworkConfig(c)
	getNetworkConfig(c)
}

func pushNetworkConfig(c *client.Client) {
	fmt.Println("\n[-] Generating & Pushing Configuration...")

	// 1. Get Demo Data
	vlans, vrfs, interfaces, routing, bgp, system := createDemoData()

	// 2. Generate XML
	configData, err := labnetdevice.GenerateEditConfig(vlans, vrfs, interfaces, routing, bgp, system)
	if err != nil {
		log.Printf("[-] XML generation error: %v", err)
		return
	}

	// 3. Send RPC
	rpc := fmt.Sprintf(`<edit-config xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <target><running/></target>
  %s
</edit-config>`, configData)

	reply, err := c.Exec(rpc)
	if err != nil {
		log.Printf("[-] Edit-Config Failed: %v", err)
		return
	}

	fmt.Println("[+] Edit-Config Configured Successfully!")
	fmt.Printf("    Message ID: %s\n", reply.MessageID)
}

func getNetworkConfig(c *client.Client) {
	fmt.Println("\n[-] Retrieving Configuration...")

	// Using explicit prefixes to avoid ambiguity
	rpc := `<get-config xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <source><running/></source>
  <filter type="subtree">
	<lnd:vlans xmlns:lnd="http://example.com/ns/lab-net-device"/>
	<lnd:vrfs xmlns:lnd="http://example.com/ns/lab-net-device"/>
	<lnd:interfaces xmlns:lnd="http://example.com/ns/lab-net-device"/>
	<lnd:routing xmlns:lnd="http://example.com/ns/lab-net-device"/>
	<lnd:bgp xmlns:lnd="http://example.com/ns/lab-net-device"/>
	<lnd:system xmlns:lnd="http://example.com/ns/lab-net-device"/>
  </filter>
</get-config>`

	reply, err := c.Exec(rpc)
	if err != nil {
		log.Printf("[-] Get-Config Failed: %v", err)
		return
	}
	fmt.Println("[+] Current Configuration (Raw XML):")
	// fmt.Println(reply.Data)

	// Parse and display structured data
	cfg, err := labnetdevice.ParseConfig(reply.Data)
	if err != nil {
		log.Printf("[-] Parse Config Failed: %v", err)
		return
	}

	fmt.Println("\n[+] Parsed Configuration structure:")
	if cfg.System != nil && cfg.System.Users != nil {
		fmt.Println("  Users:")
		for _, u := range cfg.System.Users.User {
			fmt.Printf("    - %s (%s) [%s]\n", u.UserId, u.ScreenName, u.Role)
		}
	}
	if cfg.Vlans != nil {
		fmt.Println("  VLANs:")
		for _, v := range cfg.Vlans.Vlan {
			fmt.Printf("    - ID: %d, Name: %s\n", v.Id, v.Name)
		}
	}

	if cfg.Vrfs != nil {
		fmt.Println("  VRFs:")
		for _, v := range cfg.Vrfs.Vrf {
			fmt.Printf("    - Name: %s, RD: %s\n", v.Name, v.Rd)
		}
	}

	if cfg.Interfaces != nil {
		fmt.Println("  Interfaces:")
		for _, i := range cfg.Interfaces.Interface {
			fmt.Printf("    - %s (Enabled: %v)\n", i.Name, safeBool(i.Enabled))
			if i.IPv4 != nil && len(i.IPv4.Address) > 0 {
				for _, addr := range i.IPv4.Address {
					fmt.Printf("      IP: %s/%d\n", addr.IP, safeUint8(addr.PrefixLength))
				}
			}
		}
	}

	if cfg.Routing != nil && cfg.Routing.StaticRoutes != nil {
		fmt.Println("  Static Routes:")
		for _, r := range cfg.Routing.StaticRoutes.Route {
			fmt.Printf("    - %s via %s (Dist: %d, VRF: %s)\n",
				r.Prefix,
				safeString(r.NextHop),
				safeUint8(r.Distance),
				r.Vrf,
			)
		}
	}

	if cfg.Bgp != nil {
		fmt.Printf("  BGP (Local AS: %d):\n", safeUint32(cfg.Bgp.LocalAs))
		for _, n := range cfg.Bgp.Neighbor {
			fmt.Printf("    - Neighbor: %s (Remote AS: %d, VRF: %s)\n",
				n.Address,
				safeUint32(n.RemoteAs),
				n.Vrf,
			)
		}
	}
}

func safeBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func safeUint8(p *uint8) uint8 {
	if p == nil {
		return 0
	}
	return *p
}

func safeUint32(p *uint32) uint32 {
	if p == nil {
		return 0
	}
	return *p
}

func safeString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
