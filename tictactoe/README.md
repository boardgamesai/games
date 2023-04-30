# Tic-Tac-Toe

## Overview
X and O take turns on a 3x3 grid, trying to get three in a row.

## Reference Links
* [Wikipedia](https://en.wikipedia.org/wiki/Tic-tac-toe)
* [BoardGameGeek](https://boardgamegeek.com/boardgame/11901/tic-tac-toe)

## House Rules
* X moves first

## Specifications
Your AI must implement the [`tictactoeAI`](ai/driver/tictactoe_ai.go) interface. It will be passed a [`State`](ai/driver/state.go) struct, which contains a [`Board`](board.go) representing the current state of the game.

[`Board`](board.go) is a two-dimensional `string` slice with three columns and three rows. Its values are `"X"`, `"O"`, or empty string.

The board is represented as:
```
        |        |
 [0][2] | [1][2] | [2][2]
        |        |
--------------------------
        |        |
 [0][1] | [1][1] | [2][1]
        |        |
--------------------------
        |        |
 [0][0] | [1][0] | [2][0]
        |        |
```

## Moves
Your AI must return a [`Move`](move.go) containing the `[Col][Row]` of your next move.