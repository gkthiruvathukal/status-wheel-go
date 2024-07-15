package main

import (
    "log"
    "flag"
    "net/http"
    "github.com/yourusername/myserver/pkg/status"
)

func main() {
    port := flag.String("port", "8080", "Port to run the server on")
    dbPath := flag.String("db", "./statuses.db", "Path to the SQLite3 database file")
    flag.Parse()

    db, err := status.InitDB(*dbPath)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    http.HandleFunc("/status", status.StatusHandler(db))
    http.HandleFunc("/update", status.UpdateHandler(db))
    log.Printf("Server started at :%s ", *port)
    log.Fatal(http.ListenAndServe("0.0.0.0:"+*port, nil))
}
