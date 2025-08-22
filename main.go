package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"database/sql"
	"html/template"

	_ "modernc.org/sqlite" // We have to put a _ because its a driver and will be detected as unused by compiler

	"github.com/joho/godotenv"
)

var templates *template.Template

func init() {
	fmt.Println("Initializing...")
	templates = template.Must(template.ParseGlob("static/html/*.html"))

	// Create tables in database
	db, err := sql.Open("sqlite", "database.db")
	if err != nil {
		fmt.Println("Error opening database: ", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS avatar (link TEXT)")

	if err != nil {
		fmt.Println("Error creating avatar table in database: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {

	err := templates.ExecuteTemplate(w, "index.html", nil)

	if err != nil {
		fmt.Println("Error executing template: ", err)
	}

}

func updateAvatar(link string) {

	// This will take in the avatar link from the database and check if
	// the avatar has changed using the Discord API. If it has
	// then we update it in the database.

	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		fmt.Println("Error building request: ", err)
	}

	err = godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	discord_api_token := os.Getenv("DISCORD_API_TOKEN")
	req.Header.Set("Authorization", "Bearer "+discord_api_token)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request: ", err)
	}

	type User struct {
		Id     string `json:"id"`
		Avatar string `json:"avatar"`
	}

	var user User
	json.NewDecoder(resp.Body).Decode(&user)
	avatarLink := "https://cdn.discordapp.com/avatars/" + user.Id + "/" + user.Avatar + ".png"

	fmt.Println("avatarLink: ", avatarLink+" | link: ", link)

	if avatarLink != link {
		db, err := sql.Open("sqlite", "database.db")
		if err != nil {
			fmt.Println("Error opening database: ", err)
		}
		fmt.Println("Setting avatar link to: ", avatarLink)
		_, err = db.Exec("INSERT OR REPLACE INTO avatar (link) VALUES (?)", avatarLink)
		if err != nil {
			fmt.Println("Error updating the avatar in database: ", err)
		}

		err = db.Close()

		if err != nil {
			fmt.Println("Error closing database: ", err)
		}
	}

}

func apiAvatar(w http.ResponseWriter, r *http.Request) {

	// We need to return the link from the database, but then pass it off to
	// updateAvatar to see if its changed

	// 1. Fetch link from database
	// 2. Return it

	var recheck bool
	recheck = true

	db, err := sql.Open("sqlite", "database.db")

	if err != nil {
		fmt.Println("Error opening database: ", err)
	}

	// Select link from database
	rows, err := db.Query("SELECT link FROM avatar")

	if err != nil {
		fmt.Println("Error fetching avatar link from database: ", err)
	}

	var avatarFromDB string
	rows.Next()
	err = rows.Scan(&avatarFromDB)
	if err != nil {
		fmt.Println("Error scanning rows: ", err)
	}

	if avatarFromDB == "" {
		fmt.Println("avatarFromDB is NULL")
		// Update the avatar to fill the database with a valid URL
		// None as the string so that we always generate a new URL
		updateAvatar("None")

		// Refetch the avatar from the database, should be a valid URL.
		rows, err = db.Query("SELECT link FROM avatar")

		if err != nil {
			fmt.Println("Error fetching avatar link from database: ", err)
		}

		recheck = false
	}

	rows.Next()
	err = rows.Scan(&avatarFromDB)

	if err != nil {
		fmt.Println("Error scanning rows: ", err)
	}

	fmt.Println("avatarFromDB: ", avatarFromDB)

	// We use goroutine here so that we don't have to wait for completion
	// Update the avatar using the link from the database to see if it has changed
	if recheck {
		go updateAvatar(avatarFromDB)
	}

	fmt.Fprintf(w, "%s", avatarFromDB)

}

func main() {

	http.HandleFunc("/", index)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/api/avatar", apiAvatar)

	address := "127.0.0.1:8000"

	fmt.Println("Starting server on " + address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println("Error starting server: " + err.Error())
	}

}
