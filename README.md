# boardgamesai/games

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/boardgamesai/games/tree/main.svg?style=shield)](https://dl.circleci.com/status-badge/redirect/gh/boardgamesai/games/tree/main)

This is the games library that powers [boardgames.ai](https://boardgames.ai). You can clone it to develop AIs locally, and then upload them to face off against other AIs.

---

The main point of entry for this library is [play.go](play.go).

## General usage:
```
go run play.go <game> path/to/ai_1.go path/to/ai_2.go ...
```

## Command-line options:
```
-n N     : Play N games (defaults to 1)
--random : Play a game using the sample random AIs
--raw    : Display raw JSON event data
--print  : Print the final game board/state
```

## Supported Games
1. [Four-in-a-Row](fourinarow)
1. [Game of the Amazons](amazons)
1. [Hearts](hearts)
1. [Liar's Dice](liarsdice)
1. [Reversi](reversi)
1. [Tic-Tac-Toe](tictactoe)
1. [Ultimate Tic-Tac-Toe](ulttictactoe)

## Requirements
* [Go](https://go.dev) 1.19 or higher

## Clone the repo
```
git clone git@github.com:boardgamesai/games.git
```

## Play a game
```
cd games
go run play.go --random tictactoe
```

## Develop your own AI
```
1. cp games/tictactoe/ai/example/random/random.go ~/my_ai.go
2. [ edit ~/my_ai.go ]
3. go run play.go tictactoe ~/my_ai.go games/tictactoe/ai/example/random/random.go
4. Repeat steps 2-3
```

## Constraints
1. Your AI code cannot use the network or filesystem.
1. If your AI causes a panic, it is disqualified and loses the match.
1. If your AI commits an illegal move, it is disqualified and loses the match.
1. If your AI takes longer than 15 seconds to respond with a move, it is disqualified and loses the match.
1. Your AI's code must fit in one `.go` file no larger than 1 MB.

## Notes
1. Turn order (if applicable) is always randomized
1. A disqualification is treated as a loss, for ELO calculation purposes

## Feedback
Comments / bug reports / ideas welcome at ross@boardgames.ai.