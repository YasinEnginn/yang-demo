package labnetdevice

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

const (
	Namespace   = "http://example.com/ns/lab-net-device"
	NetconfBase = "urn:ietf:params:xml:ns:netconf:base:1.0"
)

// Config represents the top-level structure for edit-config
type Config struct {
	XMLName    xml.Name    `xml:"config"`
	Xmlns      string      `xml:"xmlns,attr,omitempty"`
	Vlans      *Vlans      `xml:"vlans,omitempty"`
	Vrfs       *Vrfs       `xml:"vrfs,omitempty"`
	Interfaces *Interfaces `xml:"interfaces,omitempty"`
	Routing    *Routing    `xml:"routing,omitempty"`
	Bgp        *Bgp        `xml:"bgp,omitempty"`
	System     *System     `xml:"system,omitempty"`
}

// System Container
type System struct {
	Xmlns string `xml:"xmlns,attr,omitempty"`
	Users *Users `xml:"users,omitempty"`
}

type Users struct {
	User []User `xml:"user"`
}

type User struct {
	UserId     string `xml:"user-id"`
	ScreenName string `xml:"screen-name,omitempty"`
	Role       string `xml:"role,omitempty"` // admin | operator | readonly
}

// Vlans Container
type Vlans struct {
	Xmlns string `xml:"xmlns,attr,omitempty"`
	Vlan  []Vlan `xml:"vlan"`
}

type Vlan struct {
	Id   uint16 `xml:"id"`
	Name string `xml:"name,omitempty"`
}

// Vrfs Container
type Vrfs struct {
	Xmlns string `xml:"xmlns,attr,omitempty"`
	Vrf   []Vrf  `xml:"vrf"`
}

type Vrf struct {
	Name string `xml:"name"`
	Rd   string `xml:"rd,omitempty"`
}

// Interfaces Container
type Interfaces struct {
	Xmlns     string      `xml:"xmlns,attr,omitempty"`
	Interface []Interface `xml:"interface"`
}

type Interface struct {
	Name       string      `xml:"name"`
	Enabled    *bool       `xml:"enabled,omitempty"`
	Mtu        *uint16     `xml:"mtu,omitempty"`
	Vrf        string      `xml:"vrf,omitempty"`
	Switchport *Switchport `xml:"switchport,omitempty"`
	IPv4       *IPv4       `xml:"ipv4,omitempty"`
}

type Switchport struct {
	Mode       string  `xml:"mode,omitempty"`        // access | trunk
	AccessVlan *uint16 `xml:"access-vlan,omitempty"` // when mode=access
}

type IPv4 struct {
	Address []IPv4Address `xml:"address"`
}

type IPv4Address struct {
	IP           string `xml:"ip"`
	PrefixLength *uint8 `xml:"prefix-length,omitempty"`
}

// Routing Container
type Routing struct {
	Xmlns        string        `xml:"xmlns,attr,omitempty"`
	StaticRoutes *StaticRoutes `xml:"static-routes,omitempty"`
}

type StaticRoutes struct {
	Route []StaticRoute `xml:"route"`
}

type StaticRoute struct {
	Prefix   string `xml:"prefix"`
	Vrf      string `xml:"vrf,omitempty"`
	Distance *uint8 `xml:"distance,omitempty"`
	// Choice: next-hop-ip
	NextHop *string `xml:"next-hop,omitempty"`
	// Choice: outgoing-interface
	OutIf     *string `xml:"out-if,omitempty"`
	GatewayIP *string `xml:"gateway-ip,omitempty"`
}

// Bgp Container
type Bgp struct {
	Xmlns    string     `xml:"xmlns,attr,omitempty"`
	LocalAs  *uint32    `xml:"local-as,omitempty"`
	Neighbor []Neighbor `xml:"neighbor"`
}

type Neighbor struct {
	Address  string  `xml:"address"`
	RemoteAs *uint32 `xml:"remote-as,omitempty"`
	Vrf      string  `xml:"vrf,omitempty"`
}

// Helper to generate XML string
func (c Config) String() (string, error) {
	// We only marshal the *content* inside <config>...</config> usually,
	// but here we can marshal the whole thing and strip the outer tag if needed,
	// or return properly formatted XML.
	// For netconf edit-config, we need the content inside <config>.
	// However, usually Go XML marshaling includes the root.
	// Let's marshal internal fields individually or create a wrapper.

	// Actually, easier way:
	output, err := xml.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// GenerateEditConfig generates the content for <edit-config><target><running/></target><config>...</config></edit-config>
// It returns the inner XML structure rooted at fields.
// Since we have multiple top-level containers (vlans, vrfs, etc.), we can wrap them in a struct that doesn't print its own tag but prints children,
// OR we can just return the raw XML string constructed from parts.
//
// But to match the user request "integrate this xml", let's build a struct matches the user's provided XML structure.
// The user provided <config>...</config> block content (implied).
// Actually the user provided:
// <vlans ...> ... </vlans>
// <vrfs ...> ... </vrfs>
// ...
// inside <config>
func GenerateEditConfig(vlans *Vlans, vrfs *Vrfs, interfaces *Interfaces, routing *Routing, bgp *Bgp, system *System) (string, error) {
	// We'll create a temporary struct to marshal all together
	// We use pointers to omit empty sections
	data := struct {
		XMLName    xml.Name    `xml:"config"`
		Xmlns      string      `xml:"xmlns,attr,omitempty"` // module namespace
		Vlans      *Vlans      `xml:"vlans,omitempty"`
		Vrfs       *Vrfs       `xml:"vrfs,omitempty"`
		Interfaces *Interfaces `xml:"interfaces,omitempty"`
		Routing    *Routing    `xml:"routing,omitempty"`
		Bgp        *Bgp        `xml:"bgp,omitempty"`
		System     *System     `xml:"system,omitempty"`
	}{
		// Leave <config> without a namespace; child containers carry model NS.
		Xmlns:      "",
		Vlans:      vlans,
		Vrfs:       vrfs,
		Interfaces: interfaces,
		Routing:    routing,
		Bgp:        bgp,
		System:     system,
	}
	// The user's XML has specific namespaces on children.
	// Go's XML marshaler handles namespaces if we set them on struct fields.
	// We need to be careful not to introduce "unexpected namespace" errors.
	// If the server expects NO namespace on the top container but namespaces on children, we do this.
	// If the server expects the module namespace on the top container, we should set it there.
	// Given the error 'An unexpected namespace is present', it usually means we are sending a namespace
	// where it shouldn't be, OR we are sending the WRONG namespace.

	// Assign module namespace on each top-level child so server can resolve the YANG module.
	if vlans != nil {
		vlans.Xmlns = Namespace
	}
	if vrfs != nil {
		vrfs.Xmlns = Namespace
	}
	if interfaces != nil {
		interfaces.Xmlns = Namespace
	}
	if routing != nil {
		routing.Xmlns = Namespace
	}
	if bgp != nil {
		bgp.Xmlns = Namespace
	}
	if system != nil {
		system.Xmlns = Namespace
	}

	output, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal config: %w", err)
	}
	return string(output), nil
}

// XML -> GO
// ParseConfig unmarshals specific sections from a NETCONF <data> or <config> return.
func ParseConfig(data string) (*Config, error) {
	// NETCONF replies wrap module data in <data> and usually include namespaces.
	// We strip namespaces to keep decoding simple and accept both <config> and <data> roots.
	cleanXML, err := stripNamespaces(data)
	if err != nil {
		return nil, fmt.Errorf("failed to normalize XML: %w", err)
	}

	var cfg Config

	// First try when caller already provided a <config> root.
	if err := xml.Unmarshal([]byte(cleanXML), &cfg); err == nil {
		return &cfg, nil
	}

	// Otherwise treat payload as NETCONF <data>; wrap its children in <config>.
	var wrapper struct {
		Inner []byte `xml:",innerxml"`
	}
	if err := xml.Unmarshal([]byte(cleanXML), &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	wrapped := "<config>" + strings.TrimSpace(string(wrapper.Inner)) + "</config>"
	if err := xml.Unmarshal([]byte(wrapped), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config body: %w", err)
	}

	return &cfg, nil
}

// stripNamespaces removes XML namespaces and prefixes so encoding/xml matches our simple tags.
func stripNamespaces(input string) (string, error) {
	dec := xml.NewDecoder(strings.NewReader(input))
	var buf bytes.Buffer
	enc := xml.NewEncoder(&buf)

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			t.Name.Space = ""
			attrs := make([]xml.Attr, 0, len(t.Attr))
			for _, a := range t.Attr {
				if a.Name.Space == "xmlns" || a.Name.Local == "xmlns" {
					continue // drop namespace declarations
				}
				a.Name.Space = ""
				attrs = append(attrs, a)
			}
			t.Attr = attrs
			if err := enc.EncodeToken(t); err != nil {
				return "", err
			}
		case xml.EndElement:
			t.Name.Space = ""
			if err := enc.EncodeToken(t); err != nil {
				return "", err
			}
		case xml.CharData:
			// Clean up CharData: remove lines starting with #
			clean := cleanCharData(t)
			if err := enc.EncodeToken(clean); err != nil {
				return "", err
			}
		default:
			if err := enc.EncodeToken(tok); err != nil {
				return "", err
			}
		}
	}
	if err := enc.Flush(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func cleanCharData(data xml.CharData) xml.CharData {
	s := string(data)
	if !strings.Contains(s, "#") {
		return data
	}

	var sb strings.Builder
	lines := strings.Split(s, "\n")
	first := true
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Optional: strip trailing comments like "10 # comment"
		if idx := strings.Index(line, "#"); idx >= 0 {
			line = line[:idx]
		}

		if !first {
			sb.WriteString("\n")
		}
		sb.WriteString(line)
		first = false
	}
	return xml.CharData([]byte(sb.String()))
}
