package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "os"
    "os/exec"
    "testing"
    "time"
)

func TestIntegration(t *testing.T) {
    os.Setenv("PORT", "8081")

    // remove the database file if it exists
    _ = os.Remove("users.db")

    // run the app in a separate process
    cmd := exec.Command("go", "run", ".")
    cmd.Env = os.Environ() // preserve existing env vars

    // start the app
    if err := cmd.Start(); err != nil {
        t.Fatalf("Failed to start app process: %v", err)
    }
    defer func() {
        // make sure the server is stopped
        _ = cmd.Process.Kill()
    }()

    // wait for the server to start
    time.Sleep(2 * time.Second)

    userJSON := `{
        "external_id": "a9417568-a4fd-45ab-8705-c776bed9288c",
        "name": "Martin",
        "email": "martin@whalebone.io",
        "date_of_birth": "2000-01-01T12:12:34Z"
    }`
    // test post
    resp, err := http.Post("http://localhost:8081/save", "application/json", bytes.NewBufferString(userJSON))
    if err != nil {
        t.Fatalf("POST /save request failed: %v", err)
    }
    if resp.StatusCode != http.StatusCreated {
        t.Fatalf("Expected 201 Created, got %d", resp.StatusCode)
    }

    // test get
    externalID := "a9417568-a4fd-45ab-8705-c776bed9288c"
    getURL := "http://localhost:8081/" + externalID
    resp, err = http.Get(getURL)
    if err != nil {
        t.Fatalf("GET /:external_id request failed: %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("Expected 200 OK, got %d", resp.StatusCode)
    }
    var retrievedUser map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&retrievedUser); err != nil {
        t.Fatalf("Failed to decode GET /:external_id response: %v", err)
    }
    _ = resp.Body.Close()

    // verify that the fields
    if retrievedUser["external_id"] != externalID {
        t.Errorf("external_id mismatch: got %v", retrievedUser["external_id"])
    }
    if retrievedUser["name"] != "Martin" {
        t.Errorf("name mismatch: got %v", retrievedUser["name"])
    }
    if retrievedUser["email"] != "martin@whalebone.io" {
        t.Errorf("email mismatch: got %v", retrievedUser["email"])
    }
}
