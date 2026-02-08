# YANG Lab

A simple, hands-on playground to understand **NETCONF**, **YANG**, and **Go**.

This project demonstrates how to configure a network device programmatically using a custom YANG model (`lab-device.yang`) and a Go client.

---

## Architecture

This diagram shows how **YANG** (the rules), **XML** (the data), and **NETCONF** (the transport) work together.

```mermaid
graph TD
    %% Define Styles
    classDef yang fill:#e1f5fe,stroke:#01579b,stroke-width:2px;
    classDef xml fill:#fff9c4,stroke:#fbc02d,stroke-width:2px;
    classDef netconf fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px;

    subgraph User["Your Go Application"]
        direction TB
        Action[("Action: Configure IP")]
    end

    subgraph Transport["NETCONF (Transport Layer)"]
        RPC["<rpc message-id='101'>\n ... \n</rpc>"]:::netconf
        Operation["<edit-config>\n  <target><running/></target>\n  <config>\n    ...\n  </config>\n</edit-config>"]:::netconf
    end

    subgraph Data["XML (Data Payload)"]
        XMLContent["<interfaces>\n  <interface>\n    <name>GigabitEthernet0/0</name>\n    <ipv4>\n      <address>\n        <ip>192.0.2.1</ip>\n      </address>\n    </ipv4>\n  </interface>\n</interfaces>"]:::xml
    end

    subgraph Model["YANG (The Rules)"]
        YANGDef["module lab-device {\n  container interfaces {\n    list interface {\n      leaf name { type string; }\n      container ipv4 {\n        leaf ip { type inet:ipv4-address; }\n      }\n    }\n  }\n}"]:::yang
    end

    %% Connections
    User -->|Sends| Transport
    RPC -->|Contains| Operation
    Operation -->|Wraps| Data
    Data -.->|Validated Against| Model
```

---

## Workflow Visualization

Here is the sequence of events:

```mermaid
sequenceDiagram
    participant Client as Go Client
    participant Server as Netopeer2 (Server)
    
    Note over Client, Server: SSH Connection (Port 830)
    
    Client->>Server: <edit-config> (Set Interface IP)
    Note right of Client: Sets "GigabitEthernet0/0" <br/> to 192.0.2.1/30
    
    Server-->>Client: <ok/>
    
    Client->>Server: <get-config> (Verify Data)
    Server-->>Client: Returns Configuration Data
```

---

## Quick Start

Follow these simple steps to get everything running.

### 1. Start the Server (Docker)
Run the standard Netopeer2 container (NETCONF server).

```bash
docker run -d --name netopeer2 -p 830:830 sysrepo/sysrepo-netopeer2:latest
```

### 2. Install the YANG Model
Upload and install our custom device model (`lab-device.yang`) into the server.

```bash
# 1. Copy the YANG file to the container
docker cp yang/lab-device.yang netopeer2:/tmp/

# 2. Install the model using sysrepoctl
docker exec -it netopeer2 sysrepoctl -i /tmp/lab-device.yang
```

*(Optional)* If you face permission errors, you might need to disable NACM (NETCONF Access Control Model) inside the container.

### 3. Run the Go Client
Now, run the client to push configuration and verify it.

```bash
# Install dependencies
go mod tidy

# Run the client
go run cmd/main.go
```

---

## Project Structure

- **`yang/lab-device.yang`**: The custom data model defining our device's interface.
- **`cmd/main.go`**: The Go application that acts as the NETCONF client.

---

### Enjoy your automation!