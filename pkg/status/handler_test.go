package status

import (
    "database/sql"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    _ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("Failed to open database: %v", err)
    }

    _, err = db.Exec(`
        CREATE TABLE statuses (
            session_id TEXT PRIMARY KEY,
            status TEXT
        )
    `)
    if err != nil {
        t.Fatalf("Failed to create table: %v", err)
    }

    return db
}

func TestStatusHandler(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    _, err := db.Exec("INSERT INTO statuses (session_id, status) VALUES (?, ?)", "test", "Available")
    if err != nil {
        t.Fatalf("Failed to insert initial data: %v", err)
    }

    req, err := http.NewRequest("GET", "/status?session_id=test", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := StatusHandler(db)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    expected := `{"status":"Available"}`
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
    }
}

func TestUpdateHandler(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    var jsonStr = []byte(`{"status":"Busy"}`)
    req, err := http.NewRequest("POST", "/update?session_id=test", strings.NewReader(string(jsonStr)))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()
    handler := UpdateHandler(db)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    var status string
    err = db.QueryRow("SELECT status FROM statuses WHERE session_id = ?", "test").Scan(&status)
    if err != nil {
        t.Fatalf("Failed to query status: %v", err)
    }

    if status != "Busy" {
        t.Errorf("Expected status to be 'Busy', got '%s'", status)
    }
}
