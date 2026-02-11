package labnetdevice

import (
	"strings"
	"testing"
)

func TestGenerateEditConfig_PurposeNamespace(t *testing.T) {
	interfaces := &Interfaces{
		Interface: []Interface{
			{
				Name:    "GigabitEthernet0/0",
				Purpose: &Purpose{Value: "lndi:uplink"},
			},
		},
	}

	out, err := GenerateEditConfig(nil, nil, nil, interfaces, nil, nil, nil)
	if err != nil {
		t.Fatalf("GenerateEditConfig error: %v", err)
	}

	if !strings.Contains(out, NamespaceIdentities) {
		t.Fatalf("expected identities namespace in output, got: %s", out)
	}
	if !strings.Contains(out, NamespacePurpose) {
		t.Fatalf("expected purpose namespace in output, got: %s", out)
	}
}

func TestParseConfig_OperState(t *testing.T) {
	xmlData := `
<data>
  <interfaces xmlns="http://example.com/ns/lab-net-device">
    <interface>
      <name>GigabitEthernet1/1</name>
      <oper-status xmlns="http://example.com/ns/lab-net-device-operstate">down</oper-status>
      <last-change xmlns="http://example.com/ns/lab-net-device-operstate">2026-02-11T12:00:00Z</last-change>
      <phys-address xmlns="http://example.com/ns/lab-net-device-operstate">aa:bb:cc:dd:ee:ff</phys-address>
      <speed-mbps xmlns="http://example.com/ns/lab-net-device-operstate">1000</speed-mbps>
      <hardware-present xmlns="http://example.com/ns/lab-net-device-operstate">false</hardware-present>
      <counters xmlns="http://example.com/ns/lab-net-device-operstate">
        <in-octets>123</in-octets>
        <out-octets>456</out-octets>
      </counters>
    </interface>
  </interfaces>
</data>`

	cfg, err := ParseConfig(xmlData)
	if err != nil {
		t.Fatalf("ParseConfig error: %v", err)
	}
	if cfg.Interfaces == nil || len(cfg.Interfaces.Interface) != 1 {
		t.Fatalf("expected one interface, got: %+v", cfg.Interfaces)
	}

	iface := cfg.Interfaces.Interface[0]
	if iface.OperStatus != "down" {
		t.Fatalf("expected oper-status=down, got: %s", iface.OperStatus)
	}
	if iface.LastChange != "2026-02-11T12:00:00Z" {
		t.Fatalf("expected last-change, got: %s", iface.LastChange)
	}
	if iface.PhysAddress != "aa:bb:cc:dd:ee:ff" {
		t.Fatalf("expected phys-address, got: %s", iface.PhysAddress)
	}
	if iface.SpeedMbps == nil || *iface.SpeedMbps != 1000 {
		t.Fatalf("expected speed-mbps=1000, got: %v", iface.SpeedMbps)
	}
	if iface.HardwarePresent == nil || *iface.HardwarePresent != false {
		t.Fatalf("expected hardware-present=false, got: %v", iface.HardwarePresent)
	}
	if iface.Counters == nil || iface.Counters.InOctets == nil || iface.Counters.OutOctets == nil {
		t.Fatalf("expected counters, got: %+v", iface.Counters)
	}
	if *iface.Counters.InOctets != 123 || *iface.Counters.OutOctets != 456 {
		t.Fatalf("unexpected counters: in=%d out=%d", *iface.Counters.InOctets, *iface.Counters.OutOctets)
	}
}
