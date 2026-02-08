package client

import (
	"fmt"

	"github.com/Juniper/go-netconf/netconf"
)

// Client wraps the netconf.Session
type Client struct {
	Session *netconf.Session
}

// New creates a new NETCONF session
func New(host, user, password string) (*Client, error) {
	sshConfig := netconf.SSHConfigPassword(user, password)
	// You might want to add timeout or other config here
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
	reply, err := c.Session.Exec(netconf.RawMethod(rpc))
	if err != nil {
		return nil, err
	}
	if len(reply.Errors) > 0 {
		return nil, fmt.Errorf("RPC errors: %v", reply.Errors)
	}
	return reply, nil
}
