package main

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

// const lichessOAuthURL = "https://lichess.org/oauth"
// const lichessTokenURL = "https://lichess.org/api/token"

// var conf = &oauth2.Config{
// 	ClientID:     "chess_tui",
// 	ClientSecret: "",
// 	RedirectURL:  "http://localhost:8080/callback",
// 	Scopes:       []string{"board:play"},
// 	Endpoint: oauth2.Endpoint{
// 		AuthURL:  lichessOAuthURL,
// 		TokenURL: lichessTokenURL,
// 	},
// }

type TokenData struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64    `json:"expires_in"`
}

func main() {
	//ctx := context.Background()
	if tokensExist() {
		tok := loadTokens()
		
		fmt.Printf("%s\n%s\n%d\n", tok.AccessToken, tok.TokenType, tok.ExpiresIn)
		//client := conf.Client(ctx, tok)
		//TODO: actuall chess stuff with loaded token in client and ctx
	} else {
		// Persist for next run
		tok := loadTokens()
		fmt.Println(tok)
		// client := conf.Client(ctx, tok)
		//TODO: actuall chess stuff with loaded token in client and ctx
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

// // they don't have an oauth token stored, need to go grab it with temporary callback server
// func getOAuthToken() {

// 	// startTempServer()     // Goroutine with callback handler
// 	// openBrowser()  // Send user to Lichess
// 	// waitForCallback()     // Block until token received
// 	// saveTokensToFile(tok)
// 	verifier := oauth2.GenerateVerifier()
// 	state := rand.Text()

// 	// channels communicate values between go routines
// 	serverReadyChan := make(chan bool)  // indicates when server is ready
// 	tokChan := make(chan *oauth2.Token) // gets the token

// 	// Redirect user to consent page to ask for permission
// 	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.S256ChallengeOption(verifier))

// 	// go routine calls a function to run asynchronously
// 	go func() {
// 		// open temp web server
// 	}()

// 	// waits for server to be ready
// 	<-serverReadyChan
// 	fmt.Println("Temporary server started on localhost:8080")

// 	fmt.Println("Opening browser for authentication...")
// 	openBrowser(url)
// 	fmt.Println("Waiting for user to complete authentication...")

// 	// wait for the token to be authenticated
// 	tok := <-tokChan
// 	fmt.Println("OAuth complete! Token received")

// 	//RECEIVED TOKEN, do the chess TUI stufff
// 	fmt.Println(tok)

// 	return
// }

// // after user authorizes the application to use their account
// func callbackHandler(w http.ResponseWriter, r *http.Request) {
// 	code := r.URL.Query().Get("code")

// 	if code == "" {
// 		fmt.Fprint(w, "Authorization was not successful")
// 		os.Exit(1)
// 	}

// 	// writes to the http response so grab that and use for the oauth token
// 	fmt.Fprint(w, "Authorization was successful! You can close this window and get back to your terminal")
// 	// fmt.Printf("Received code: %s, state: %s\n", code, state)

// 	fmt.Println("ready to call the other endpoint for access/oauth token")

// 	// send authroization code to lichesTokenURL get access token
// 	tok, err := conf.Exchange(r.Context(), code, oauth2.VerifierOption(verifier))

// 	// Save the entire token source state
// 	tokenSource := conf.TokenSource(r.Context(), tok)
// 	persistedToken, _ := tokenSource.Token() // Get current token
// 	saveToFile("tokens.json", persistedToken)

// 	if err != nil {
// 		panic(err)
// 	}

// }

// // https://www.ziye.dev/posts/go-server-for-oauth-callbacks/
// func openBrowser(url string) {
// 	var err error

// 	switch runtime.GOOS {
// 	case "linux":
// 		err = exec.Command("xdg-open", url).Start()
// 	case "windows":
// 		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
// 	case "darwin":
// 		err = exec.Command("open", url).Start()
// 	default:
// 		err = fmt.Errorf("unsupported platform")
// 	}
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// // saves the token to the json file
// func saveToFile(filename string, token *oauth2.Token) {
// 	return
// }
