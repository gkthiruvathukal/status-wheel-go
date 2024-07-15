package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "bytes"
    "os"
    "time"
)

// Struct to parse JSON responses
type StatusResponse struct {
    Status string `json:"status"`
}

func main() {
    // Define the subcommands
    checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
    setCmd := flag.NewFlagSet("set", flag.ExitOnError)
    monitorCmd := flag.NewFlagSet("monitor", flag.ExitOnError)

    // Common flags
    sessionID := checkCmd.String("session_id", "", "Session ID")
    serverAddr := checkCmd.String("server", "http://localhost:8080", "Server address")

    setSessionID := setCmd.String("session_id", "", "Session ID")
    setStatus := setCmd.String("status", "", "Status to set")
    setServerAddr := setCmd.String("server", "http://localhost:8080", "Server address")

    monitorSessionID := monitorCmd.String("session_id", "", "Session ID")
    monitorInterval := monitorCmd.Int("interval", 10, "Polling interval in seconds")
    monitorServerAddr := monitorCmd.String("server", "http://localhost:8080", "Server address")

    // Parse the subcommands
    if len(os.Args) < 2 {
        fmt.Println("Expected 'check', 'set', or 'monitor' subcommands")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "check":
        checkCmd.Parse(os.Args[2:])
        if *sessionID == "" {
            log.Fatal("session_id is required")
        }
        checkStatus(*serverAddr, *sessionID)
    case "set":
        setCmd.Parse(os.Args[2:])
        if *setSessionID == "" {
            log.Fatal("session_id is required")
        }
        if *setStatus == "" {
            log.Fatal("status is required")
        }
        setStatusFunc(*setServerAddr, *setSessionID, *setStatus)
    case "monitor":
        monitorCmd.Parse(os.Args[2:])
        if *monitorSessionID == "" {
            log.Fatal("session_id is required")
        }
        monitorStatus(*monitorServerAddr, *monitorSessionID, *monitorInterval)
    default:
        fmt.Println("Expected 'check', 'set', or 'monitor' subcommands")
        os.Exit(1)
    }
}

func checkStatus(serverAddr, sessionID string) {
    resp, err := http.Get(fmt.Sprintf("%s/status?session_id=%s", serverAddr, sessionID))
    if err != nil {
        log.Fatalf("Failed to check status: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Failed to read response: %v", err)
    }

    var statusResp StatusResponse
    if err := json.Unmarshal(body, &statusResp); err != nil {
        log.Fatalf("Failed to parse response: %v", err)
    }

    fmt.Printf("Status for session %s: %s ", sessionID, statusResp.Status)
}

func setStatusFunc(serverAddr, sessionID, status string) {
    data := map[string]string{"status": status}
    jsonData, err := json.Marshal(data)
    if err != nil {
        log.Fatalf("Failed to marshal JSON: %v", err)
    }

    resp, err := http.Post(fmt.Sprintf("%s/update?session_id=%s", serverAddr, sessionID), "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Fatalf("Failed to set status: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Fatalf("Failed to set status, server responded with status: %s", resp.Status)
    }

    fmt.Printf("Status for session %s set to %s ", sessionID, status)
}

func monitorStatus(serverAddr, sessionID string, interval int) {
    ticker := time.NewTicker(time.Duration(interval) * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            checkStatus(serverAddr, sessionID)
        }
    }
}
