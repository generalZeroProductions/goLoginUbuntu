package main

import (
   
    "fmt"
    "log"
    "net/http"
    "os/exec"
    "github.com/gorilla/websocket"
	
	"time"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)
// var dotBroadcast = make(chan bool)

var timer *time.Timer
var done = make(chan bool)

func main() {
    http.HandleFunc("/ws", handleConnections)
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    http.HandleFunc("/", loginHandler)
    http.HandleFunc("/newSite", newSiteHandler)
	 http.HandleFunc("/viewSites", viewSitesHandler)
	  http.HandleFunc("/deleteSite", deleteSiteHandler)
    go handleMessages() // Start the message handling goroutine


    fmt.Println("Server is running on port 2000...")
    log.Fatal(http.ListenAndServe(":2000", nil))
}


func newSiteHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        websiteName := r.FormValue("websitename")
        if websiteName == "" {
            fmt.Fprintf(w, "Website name cannot be empty")
            return
        }

        cmd := exec.Command("/bin/bash", "newSite.sh", websiteName, username, password)
         feedback := make(chan string)
        go submitHandler(cmd, feedback)
        
        message := <-feedback
        renderNewSitePage(w, message)
    } else {
       
		renderNewSitePage(w, "")
    }
}


func deleteSiteHandler(w http.ResponseWriter, r *http.Request) {
     if r.Method == http.MethodPost {
        websiteName := r.FormValue("websitename")
        if websiteName == "" {
            fmt.Fprintf(w, "Website name cannot be empty")
            return
        }
        cmd := exec.Command("/bin/bash", "deleteSite.sh", websiteName, username, password)
         feedback := make(chan string)
        go submitHandler(cmd, feedback)
        
        message := <-feedback
        renderNewSitePage(w, message)
    } else {
		renderNewSitePage(w, "")
    }
}

func viewSitesHandler(w http.ResponseWriter, r *http.Request) {
	logToFile("vie handler here")
     if r.Method == http.MethodPost {
        cmd := exec.Command("/bin/bash", "viewSites.sh")
         feedback := make(chan string)
        go submitHandler(cmd, feedback)
        
        message := <-feedback
        renderNewSitePage(w, message)
    } else {
		renderNewSitePage(w, "")
    }
}












