package client

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/Juniper/go-netconf/netconf"
)

// Client wraps the netconf.Session
type Client struct {
	Session *netconf.Session
}

// New creates a new NETCONF session
func New(host, user, password string) (*Client, error) {
	if strings.TrimSpace(host) == "" {
		return nil, fmt.Errorf("host is required")
	}
	if strings.TrimSpace(user) == "" {
		return nil, fmt.Errorf("user is required")
	}
	if password == "" {
		return nil, fmt.Errorf("password is required")
	}
	host = ensurePort(host, "830")
	sshConfig := netconf.SSHConfigPassword(user, password)
	sshConfig.Timeout = 10 * time.Second
	session, err := netconf.DialSSH(host, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", host, err)
	}
	return &Client{Session: session}, nil
}

// Close closes the session
func (c *Client) Close() {
	if c.Session != nil {
		c.Session.Close()
	}
}

// Exec executes a raw RPC method
func (c *Client) Exec(rpc string) (*netconf.RPCReply, error) {
	if c.Session == nil {
		return nil, fmt.Errorf("netconf session is nil")
	}
	trim := strings.TrimSpace(rpc)
	if strings.HasPrefix(trim, "<rpc") || strings.HasPrefix(trim, "<rpc ") {
		return nil, fmt.Errorf("rpc must not include <rpc> wrapper")
	}
	reply, err := c.Session.Exec(netconf.RawMethod(rpc))
	if err != nil {
		return nil, fmt.Errorf("netconf exec failed: %w", err)
	}
	if len(reply.Errors) > 0 {
		return reply, fmt.Errorf("RPC errors:\n%s", formatRPCErrors(reply.Errors))
	}
	return reply, nil
}

func ensurePort(host, port string) string {
	if host == "" {
		return host
	}
	host = strings.TrimSuffix(strings.TrimSpace(host), ":")
	if _, _, err := net.SplitHostPort(host); err == nil {
		return host
	}
	if ip := net.ParseIP(host); ip != nil && ip.To4() == nil {
		return "[" + host + "]:" + port
	}
	return host + ":" + port
}

func formatRPCErrors(errs []netconf.RPCError) string {
	if len(errs) == 0 {
		return ""
	}
	lines := make([]string, 0, len(errs))
	for i, e := range errs {
		lines = append(lines, fmt.Sprintf(
			"%d) type=%s tag=%s severity=%s message=%q path=%q",
			i+1,
			getRPCFieldString(e, "Type", "ErrorType"),
			getRPCFieldString(e, "Tag", "ErrorTag"),
			getRPCFieldString(e, "Severity", "ErrorSeverity"),
			getRPCFieldString(e, "Message", "ErrorMessage"),
			getRPCFieldString(e, "Path", "ErrorPath"),
		))
	}
	return strings.Join(lines, "\n")
}

func getRPCFieldString(e netconf.RPCError, names ...string) string {
	v := reflect.ValueOf(e)
	for _, name := range names {
		f := v.FieldByName(name)
		if f.IsValid() && f.Kind() == reflect.String {
			return f.String()
		}
	}
	return ""
}
