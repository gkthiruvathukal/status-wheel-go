package status

import (
    "database/sql"
    "encoding/json"
    "net/http"
    _ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS statuses (
            session_id TEXT PRIMARY KEY,
            status TEXT
        )
    `)
    if err != nil {
        return nil, err
    }

    return db, nil
}

func StatusHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        sessionID := r.URL.Query().Get("session_id")
        if sessionID == "" {
            http.Error(w, "Missing session_id", http.StatusBadRequest)
            return
        }

        var status string
        err := db.QueryRow("SELECT status FROM statuses WHERE session_id = ?", sessionID).Scan(&status)
        if err != nil {
            if err == sql.ErrNoRows {
                status = "Unknown"
            } else {
                http.Error(w, "Database error", http.StatusInternalServerError)
                return
            }
        }

        json.NewEncoder(w).Encode(map[string]string{"status": status})
    }
}

func UpdateHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
            return
        }

        sessionID := r.URL.Query().Get("session_id")
        if sessionID == "" {
            http.Error(w, "Missing session_id", http.StatusBadRequest)
            return
        }

        var data map[string]string
        if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        status, ok := data["status"]
        if !ok {
            http.Error(w, "Missing status", http.StatusBadRequest)
            return
        }

        _, err := db.Exec("INSERT OR REPLACE INTO statuses (session_id, status) VALUES (?, ?)", sessionID, status)
        if err != nil {
            http.Error(w, "Database error", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusOK)
    }
}
