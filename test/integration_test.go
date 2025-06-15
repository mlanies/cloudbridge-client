package test

import (
    "os/exec"
    "testing"
    "time"
    "net"
)

func TestRelayIntegration(t *testing.T) {
    relay := exec.Command("./relay-server", "--debug")
    if err := relay.Start(); err != nil {
        t.Fatalf("Не удалось запустить relay-server: %v", err)
    }
    defer relay.Process.Kill()
    time.Sleep(2 * time.Second)

    client := exec.Command("./cloudbridge-client", "--config", "./testdata/config-test.yaml")
    if err := client.Start(); err != nil {
        t.Fatalf("Не удалось запустить cloudbridge-client: %v", err)
    }
    defer client.Process.Kill()
    time.Sleep(2 * time.Second)

    conn, err := net.DialTimeout("tcp", "localhost:3389", 2*time.Second)
    if err != nil {
        t.Fatalf("Туннель не поднят: %v", err)
    }
    conn.Close()
} 