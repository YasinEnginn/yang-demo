package main

import (
	"fmt"
	"log"

	"github.com/Juniper/go-netconf/netconf"
)

func main() {
	s, err := netconf.DialSSH(
		"127.0.0.1:830",
		netconf.SSHConfigPassword("netconf", "netconf"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// 1) WRITE (edit-config)
	editRPC := `<edit-config xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <target><running/></target>
  <config>
    <interfaces xmlns="http://example.com/ns/lab-device">
      <interface>
        <name>GigabitEthernet0/0</name>
        <enabled>true</enabled>
        <description>WAN uplink</description>
        <mtu>1500</mtu>
        <ipv4>
          <address>
            <ip>192.0.2.1</ip>
            <prefix-length>30</prefix-length>
          </address>
        </ipv4>
      </interface>
    </interfaces>
  </config>
</edit-config>`

	editReply, err := s.Exec(netconf.RawMethod(editRPC))
	if err != nil {
		log.Fatalf("edit-config failed: %v", err)
	}
	fmt.Println("=== EDIT-CONFIG DATA ===\n", editReply.Data)
	fmt.Println("=== EDIT-CONFIG RAW  ===\n", editReply.RawReply)

	// 2) READ (get-config) - verify like "show running-config"
	getRPC := `<get-config xmlns="urn:ietf:params:xml:ns:netconf:base:1.0">
  <source><running/></source>
  <filter type="subtree">
    <interfaces xmlns="http://example.com/ns/lab-device"/>
  </filter>
</get-config>`

	getReply, err := s.Exec(netconf.RawMethod(getRPC))
	if err != nil {
		log.Fatalf("get-config failed: %v", err)
	}
	fmt.Println("=== GET-CONFIG DATA ===\n", getReply.Data)
	fmt.Println("=== GET-CONFIG RAW  ===\n", getReply.RawReply)
}
