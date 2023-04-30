# Ultimate Tic-Tac-Toe

## Overview
A variation on tic-tac-toe where the goal is to win a grid comprised of smaller subgrids.

## Reference Links
* [Wikipedia](https://en.wikipedia.org/wiki/Ultimate_tic-tac-toe)
* [BoardGameGeek](https://boardgamegeek.com/boardgame/9898/tic-tac-toe-times-10)

## House Rules
* X moves first, and can play anywhere on the first move
* If you are sent to a subgrid that already has a winner, you can move anywhere

## Specifications
Your AI must implement the [`ulttictactoeAI`](ai/driver/ulttictactoe_ai.go) interface. It will be passed a [`State`](ai/driver/state.go) struct, which contains a [`Board`](board.go) representing the current state of the game. Note: Ultimate Tic-Tac-Toe borrows data structures from this library's [Tic-Tac-Toe](../tictactoe) implementation.

[`Board`](board.go) is a struct with the following fields:
* `SubGrids` - a 3x3 slice of type [`tictactoe.Board`](../tictactoe/board.go) that represents the board (see [Tic-Tac-Toe](../tictactoe) for its usage). `SubGrids` is represented as:

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
* `Grid` - a [`tictactoe.Board`](../tictactoe/board.go) that tracks the winners of the subgrids (if any)
* `NextPlay` - the `Coords` (`[Col,Row]`) of the subgrid for the next play, if `nil` then the next play can be anywhere

## Moves
Your AI must return a [`Move`](move.go) containing the `[Col][Row]` of the subgrid in which to play, and the `[SubCol][SubRow]` to play in that grid. Note: if `NextPlay` is not obeyed, your AI will be disqualified.