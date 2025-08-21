package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func index(w http.ResponseWriter, r *http.Request) {

	file, err := os.Open("static/html/index.html")
	if err != nil {
		fmt.Println("Error opening HTML file: ", err)
	}
	_, err = io.Copy(w, file)
	if err != nil {
		fmt.Println("Error copying file to http response writer: ", err)
	}

}

func apiAvatar(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		fmt.Println("Error building request: ", err)
	}

	err = godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file: ", err)
	}

	discord_api_token := os.Getenv("DISCORD_API_TOKEN")
	fmt.Println(discord_api_token)
	req.Header.Set("Authorization", "Bearer "+discord_api_token)
	fmt.Println(req)

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

	fmt.Fprintf(w, "%s", "https://cdn.discordapp.com/avatars/"+user.Id+"/"+user.Avatar+".png")

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
