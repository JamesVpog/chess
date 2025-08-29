package main

import (
	"crypto/rand"
	"fmt"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
	
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

const lichessOAuthURL = "https://lichess.org/oauth"
const lichessTokenURL = "https://lichess.org/api/token"

var conf = &oauth2.Config{
		ClientID:     "chess_tui",
		ClientSecret: "",
		RedirectURL: "http://localhost:8080/callback",
		Scopes:       []string{"board:play"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  lichessOAuthURL,
			TokenURL: lichessTokenURL,
		},
	}



func main() {

	if os.Getenv("SESSION_KEY") == "" {
		fmt.Println("Please set up a secure environment variable called SESSION_KEY. See the README for more details")
		os.Exit(1)
	}

	// temp web server for OAuth2
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/callback", callbackHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

// sets up user for logging in, will redirect to actual lichess login page
// starts the OAuth2 Authorization Code Flow with PKCE
func loginHandler(w http.ResponseWriter, r *http.Request) {

	// credit: example code in the github.com/gorilla/sessions README
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := store.Get(r, "session-name")

	// Set some session values.
	verifier := oauth2.GenerateVerifier() 
	state := rand.Text()

	session.Values["verifier"] = verifier
	session.Values["state"] = state

	// Save it before we write to the response/return from the handler.
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))
	
	// Redirect the user's browser to Lichess's authorization URL 
	http.Redirect(w, r, url, http.StatusFound)

}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	
	fmt.Println("Authorization succeeded!")

	// writes to the http response so grab that and use for the oauth token
	fmt.Fprintf(w, "Authorization code: %s\nState: %s", code, state)
	// fmt.Printf("Received code: %s, state: %s\n", code, state)
	
	fmt.Println("ready to call the other endpoint for access/oauth token")
	//TODO: send authroization code to lichesTokenURL get access token 
	
	

}

//TODO: how to communicate using board api ...
//TODO: how to keep session open and listen for bot moves

//TODO: how to take over terminal with chess TUI while server runs? I need server to refresh token every now and then
//TODO; learn concurrency probably