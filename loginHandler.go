package main
import (
    "fmt"
     "html/template"
    "net/http"
    "github.com/msteinert/pam"
    "os"
    "log"
)
var username string
var password string


func loginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        username = r.FormValue("username")
        password = r.FormValue("password")

        authenticated := authenticateUser(username, password)
        if authenticated {
            http.Redirect(w, r, "/newSite", http.StatusSeeOther)
            logToFile("Authentication successful for user: " + username)
        } else {
            feedback := "Authentication failed. Please try again."
            logToFile("Authentication failed for user: " + username)
            renderLoginPage(w, feedback)
        }
    } else {
        renderLoginPage(w, "")
    }
}





func renderLoginPage(w http.ResponseWriter, feedback string) {
    tmpl, err := template.ParseFiles("login.html")
    if err != nil {
        log.Fatalf("Error parsing template: %v", err)
    }
    data := struct {
        Feedback string
    }{
        Feedback: feedback,
    }
    tmpl.Execute(w, data)
}

func renderNewSitePage(w http.ResponseWriter, feedback string) {
    tmpl, err := template.ParseFiles("newSite.html")
    if err != nil {
        log.Fatalf("Error parsing template: %v", err)
    }
    data := struct {
        Feedback string
    }{
        Feedback: feedback,
    }
    tmpl.Execute(w, data)
}



func authenticateUser(username, password string) bool {
    t, err := pam.StartFunc("login", username, func(s pam.Style, msg string) (string, error) {
        if s == pam.PromptEchoOff || s == pam.PromptEchoOn {
            return password, nil
        }
        return "", nil
    })
    if err != nil {
        logToFile("Error starting PAM transaction: " + err.Error())
        return false
    }
    if err = t.Authenticate(0); err != nil {
        logToFile("Error authenticating user: " + err.Error())
        return false
    }
    return true
}





func logToFile(message string) {
    file, err := os.OpenFile("login_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Println("Failed to open file:", err)
        return
    }
    defer file.Close()

    _, err = file.WriteString(message + "\n")
    if err != nil {
        fmt.Println("Failed to write to file:", err)
    }
}