# Chess

I like to play chess, but I don't want to leave my terminal.

SOO like an engineer, this is my attempt to build a tui (terminal user interface) to play chess.

To actually connect to Lichess, we must be secure so I want users to generate their own session secret.
1. Generate a secret key and store it in environment with the following command (you many need to install openssl):
   `export SESSION_KEY=$(openssl rand -base64 32)`
   Verify that the output of `echo $SESSION_KEY` exists

3. Run the application 
   It will os.Exit(1) if you don't have one so make sure you have a .env with a SESSION_SECRET variable


roadmap of things to learn in this repo:
- [ ] lichess.org endpoints, I think there is a way to communicate with a chess bot 
- [ ] OATUH2.0 with pcke very nice... 
- [ ] different terminal UI packages or create my own...(probably not)

roadmap of features I want:
- [ ] input moves via terminal 
- [ ] send move to the engine/player
- [ ] receive new board position and prompt user again for move
- [ ] allow user to save the game (FEN or PGN or idk how else chess games are saved)
- [ ] extra stuff:
  - [ ] chess offline (using bitboard, engine creation, etc, etc)
  - [ ] multiplayer
  - [ ] user history 
