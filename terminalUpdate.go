package main

import (
    "bufio"
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "github.com/gorilla/websocket"
)








func submitHandler(cmd *exec.Cmd, feedback chan<- string) {
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        feedback <- fmt.Sprintf("Failed to get stdout pipe: %s", err)
        return
    }
    if err := cmd.Start(); err != nil {
        feedback <- fmt.Sprintf("Failed to start command: %s", err)
        return
    }

    go func() {
        scanner := bufio.NewScanner(stdout)
        for scanner.Scan() {
            line := scanner.Text()
            logToFile(line)     // Log to file
            terminalBroadcast(line) // Send to WebSocket clients
        }
        if err := scanner.Err(); err != nil {
            log.Println("Error reading stdout:", err)
        }

        if err := cmd.Wait(); err != nil {
            log.Println("Command finished with error:", err)
            feedback <- fmt.Sprintf("Command finished with error: %s", err)
        } else {
            feedback <- "COMPLETE"
        }
    }()
}










func handleConnections(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    defer conn.Close()

    clients[conn] = true

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error: %v", err)
            delete(clients, conn)
            break
        }
    }
}


func handleMessages() {
    for {
        message := <-broadcast
        if message == "COMPLETE" {
            return // Stop the goroutine when script completes
        }
        for client := range clients {
            err := client.WriteMessage(websocket.TextMessage, []byte(message))
            if err != nil {
                delete(clients, client)
            }
        }
    }
}


func terminalBroadcast(message string) {
 
	broadcast <- message
	
}