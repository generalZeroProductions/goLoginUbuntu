package main

import (
    "bufio"
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "github.com/gorilla/websocket"
	// "strings"
	"time"
)

func submitHandler(cmd *exec.Cmd,feedback chan<- string) {
  logToFile("in submit");
    stdout, err := cmd.StdoutPipe()
    if err != nil {
          feedback <- fmt.Sprintf("Failed to get stdout pipe: %s", err)
        return
    }

    if err := cmd.Start(); err != nil {
        feedback <- fmt.Sprintf("Failed to get stdout pipe: %s", err)
        return
    }

    // Channel to signal the start of the script
    scriptDone := make(chan bool)

    go func() {
        timer := time.NewTimer(1 * time.Second)
        defer timer.Stop()

        go func() {
            for {
                select {
                case <-timer.C:
                    // Timer expired, send a dot to all clients
                    for client := range clients {
                        err := client.WriteMessage(websocket.TextMessage, []byte("."))
                        if err != nil {
                            log.Printf("WebSocket error: %v", err)
                            client.Close()
                            delete(clients, client)
                        }
                    }
                    // Reset the timer
                    timer.Reset(1 * time.Second)
                case <-scriptDone:
                    return
                }
            }
        }()

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
        }

        scriptDone <- true
        terminalBroadcast("COMPLETE")
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
    // var formattedMessage string
	broadcast <- message
	// cleanedString := strings.TrimPrefix(message, "success:")
    // if strings.Contains(message, "success") {
    //     // if message == "success: website creation complete" {
    //     //     formattedMessage = "<br><span style='color:chartreuse'>" + cleanedString + "</span><br>"
    //     // } else {
	// 	// 	if message =="success: directory created" {
	// 	// 		 formattedMessage = "<span>" + cleanedString + "</span><br>"
	// 	// 	} else {
	// 	// 		 formattedMessage = "<br> <span>" + cleanedString + "</span><br>"
	// 	// 	}
    //     // }
    //     broadcast <- message
    // } else if strings.Contains(message, "failure") {
    //     formattedMessage = "<br><span style='color:red'>" + message + "</span><br>"
    //     broadcast <- formattedMessage
    // } else if message == "." {
    //     broadcast <- message
    // } 
	
}