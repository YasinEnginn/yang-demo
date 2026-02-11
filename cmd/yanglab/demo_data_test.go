package main

import "testing"

func TestCreateDemoData_Preprov(t *testing.T) {
	_, _, _, interfaces, _, _, _ := createDemoData("default", true)
	if interfaces == nil {
		t.Fatal("interfaces is nil")
	}

	found := false
	for _, iface := range interfaces.Interface {
		if iface.Name == "GigabitEthernet1/1" {
			found = true
			break
		}
	}

	if !found {
		t.Fatal("expected pre-provisioned interface GigabitEthernet1/1")
	}

	_, _, _, interfaces, _, _, _ = createDemoData("default", false)
	for _, iface := range interfaces.Interface {
		if iface.Name == "GigabitEthernet1/1" {
			t.Fatal("did not expect pre-provisioned interface when preprov=false")
		}
	}
}
