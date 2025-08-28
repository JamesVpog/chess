package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"os"
	
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))


const lichessOAuthURL = "https://lichess.org/oauth"

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
	codeVerifier := rand.Text()
	codeChallenge := generateCodeChallenge(codeVerifier)
	state := rand.Text()

	session.Values["code_challenge"] = codeChallenge
	session.Values["state"] = state

	// Save it before we write to the response/return from the handler.
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build the lichessOAuthURL with the correct parameters
	req, err := http.NewRequest(http.MethodGet, lichessOAuthURL, nil)
	if err != nil {
	fmt.Printf("client: could not create request: %s\n", err)
	os.Exit(1)
	}

	q := req.URL.Query() // Get a copy of the query values.
	// add values to the set
	q.Add("response_type", "code") 
	q.Add("client_id", "chess_tui") 
	q.Add("redirect_uri", "http://localhost:8080/callback") 
	q.Add("code_challenge_method", "S256")
	q.Add("code_challenge",codeChallenge)  
	q.Add("state", state)

	req.URL.RawQuery = q.Encode() // Encode and assign back to the original query.
 
	// Redirect the user's browser to Lichess's authorization URL 
	http.Redirect(w, r, req.URL.String(), http.StatusFound)

}

//returns a codeChallenge given a codeVerifier
func generateCodeChallenge(codeVerifier string) (codeChallenge string) {
	 
	h := sha256.New() //create sha256 hash

	h.Write([]byte(codeVerifier)) // hash it  

	// hash it? idk what h.Sum(nil) actually does  and base64encode it
	return base64.URLEncoding.EncodeToString(h.Sum(nil)) 

}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	
	fmt.Println("Authorization succeeded!")

	// writes to the http response so grab that and use for the oauth token
	fmt.Fprintf(w, "Authorization code: %s\nState: %s", code, state)
	// fmt.Printf("Received code: %s, state: %s\n", code, state)
	
	fmt.Println("ready to call the other endpoint for access/oauth token")
	//TODO: send oauth token and get access to everything as a user
}

//TODO: how to communicate using board api ...
//TODO: how to keep session open and listen for bot moves

//TODO: how to take over terminal with chess TUI while server runs? I need server to refresh token every now and then
//TODO; learn concurrency probably