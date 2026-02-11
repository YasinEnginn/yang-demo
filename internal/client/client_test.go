package client

import (
	"strings"
	"testing"

	"github.com/Juniper/go-netconf/netconf"
)

func TestNewValidation(t *testing.T) {
	if _, err := New("", "user", "pass"); err == nil {
		t.Fatal("expected error for empty host")
	}
	if _, err := New("127.0.0.1", "", "pass"); err == nil {
		t.Fatal("expected error for empty user")
	}
	if _, err := New("127.0.0.1", "user", ""); err == nil {
		t.Fatal("expected error for empty password")
	}
}

func TestEnsurePort(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"127.0.0.1", "127.0.0.1:830"},
		{"127.0.0.1:830", "127.0.0.1:830"},
		{"example.com", "example.com:830"},
		{"example.com:2022", "example.com:2022"},
		{"example.com:", "example.com:830"},
		{"2001:db8::1", "[2001:db8::1]:830"},
	}

	for _, tt := range tests {
		if got := ensurePort(tt.in, "830"); got != tt.want {
			t.Fatalf("ensurePort(%q)=%q want=%q", tt.in, got, tt.want)
		}
	}
}

func TestExecNilSession(t *testing.T) {
	c := &Client{}
	if _, err := c.Exec("<get-config/>"); err == nil {
		t.Fatal("expected error for nil session")
	}
}

func TestExecRejectsRPCWrapper(t *testing.T) {
	c := &Client{Session: &netconf.Session{}}
	_, err := c.Exec("<rpc><get-config/></rpc>")
	if err == nil || !strings.Contains(err.Error(), "rpc must not include") {
		t.Fatalf("expected wrapper error, got: %v", err)
	}
}

func TestFormatRPCErrorsEmpty(t *testing.T) {
	if got := formatRPCErrors(nil); got != "" {
		t.Fatalf("expected empty string, got: %q", got)
	}
}
