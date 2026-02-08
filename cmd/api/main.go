package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"yang/internal/client"
)

// Global client (in a real app, use dependency injection or a handler struct)
var netconfClient *client.Client

func main() {
	// 1. Connect to NETCONF Server
	var err error
	netconfClient, err = client.New("127.0.0.1:830", "netconf", "netconf")
	if err != nil {
		log.Fatalf("Failed to connect to NETCONF server: %v", err)
	}
	defer netconfClient.Close()
	log.Println("Connected to NETCONF Server")

	// 2. Define Routes
	http.HandleFunc("/api/v1/network", handleNetworkConfig)

	// 3. Start Server
	port := ":8080"
	log.Printf("Starting API server on %s...", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

// handleNetworkConfig handles GET and POST for network config
func handleNetworkConfig(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getNetwork(w, r)
	case http.MethodPost:
		updateNetwork(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getNetwork(w http.ResponseWriter, r *http.Request) {
	// Placeholder: In a real implementation, you would call netconfClient.Exec(...)
	// and parse/marshal proper structs from labnetdevice package.

	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status":  "not implemented yet",
		"message": "Use CLI for now, or implement full labnetdevice mapping here",
	}
	json.NewEncoder(w).Encode(response)
}

func updateNetwork(w http.ResponseWriter, r *http.Request) {
	// Placeholder
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "received"}`)
}
