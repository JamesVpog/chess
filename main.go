package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"golang.org/x/oauth2"
)

const lichessOAuthURL = "https://lichess.org/oauth"
const lichessTokenURL = "https://lichess.org/api/token"
const rapidTV = "https://lichess.org/api/tv/rapid/feed"

var conf = &oauth2.Config{
	ClientID:     "chess_tui",
	ClientSecret: "",
	RedirectURL:  "http://localhost:8080/callback",
	Scopes:       []string{"email:read"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  lichessOAuthURL,
		TokenURL: lichessTokenURL,
	},
}

type TokenData struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func main() {
	ctx := context.Background()
	if !tokensExist() {
		getOAuthToken()
	}
	tok := loadTokens()
	client := conf.Client(ctx, &tok)
	featuredTVPrintGame(client)

}

//TODO: how to ingest the most featured game and view in terminal? test for ingesting real game!

// so how do we represent the board in terminal?? how do we turn PEN data into board state? 
// we need 8x8 grid of 64 squares in the terminal


// lets just first figure out how to ingest nd-json continuously
func featuredTVPrintGame(client *http.Client) {
	resp, err := client.Get(rapidTV)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	// until the game ends, keep printing the chess moves
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}


// token.json holds the oauth token necessary for interacting with lichess api
// checks if token.json file exists
func tokensExist() bool {
	_, err := os.Open("token.json")
	if os.IsNotExist(err) {
		return false
	}
	if os.IsExist(err) && err != nil {
		panic(err)
	}
	return true
}

// loads token from token.json
func loadTokens() (token oauth2.Token) {
	data, err := os.ReadFile("token.json")
	if err != nil {
		panic(err)
	}
	var tokenData TokenData

	err = json.Unmarshal(data, &tokenData)
	if err != nil {
		panic(err)
	}

	token = oauth2.Token{
		AccessToken: tokenData.AccessToken,
		TokenType:   tokenData.TokenType,
		ExpiresIn:   tokenData.ExpiresIn,
	}

	return
}

// they don't have an oauth token stored, need to go grab it with temporary callback server
// save to the token.json
func getOAuthToken() {

	verifier := oauth2.GenerateVerifier()
	state := rand.Text()

	// user login URL
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

	// channels communicate values between go routines
	tokChan := make(chan *oauth2.Token) // gets the token

	// schematic of temporary web server
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		stateFromLichess := r.URL.Query().Get("state")

		if code == "" || stateFromLichess != state {
			fmt.Fprint(w, "Authorization was not successful")
			os.Exit(1)
		}

		tok, err := conf.Exchange(r.Context(), code, oauth2.VerifierOption(verifier))
		if err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}

		// Send token to main goroutine when it is received
		tokChan <- tok
		fmt.Fprintf(w, "Authentication successful! You can close this window.")
	})

	// start the server
	server := &http.Server{Addr: "localhost:8080", Handler: mux}
	go server.ListenAndServe()

	fmt.Println("Opening browser for authentication...")
	openBrowser(url)
	fmt.Println("Waiting for user to complete authentication...")

	// wait for the token to be authenticated
	token := <-tokChan

	// I have no idea what this does but it works
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	saveToFile(token)
}

// https://www.ziye.dev/posts/go-server-for-oauth-callbacks/
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// saves the token to the json file
func saveToFile(token *oauth2.Token) {
	b, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("token.json", b, 0644)

	if err != nil {
		panic(err)
	}
}
