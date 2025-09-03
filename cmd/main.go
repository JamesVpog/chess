package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"

	"github.com/JamesVPog/go-chess/api"
)


const rapidTV = "https://lichess.org/api/tv/rapid/feed"


type TokenData struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func main() {
	ctx := context.Background()
	if !api.TokensExist() {
		api.GetOAuthToken()
	}
	tok := api.LoadTokens()
	client := api.Conf.Client(ctx, &tok)
	featuredTVPrintGame(client)

	// building the UI for the chess game

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
