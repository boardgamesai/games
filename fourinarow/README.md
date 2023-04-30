# Four-in-a-Row

## Overview
Players take turns dropping discs into a 6x7 grid, trying to get four in a row.

## Reference Links
* [Wikipedia](https://en.wikipedia.org/wiki/Connect_Four)
* [BoardGameGeek](https://boardgamegeek.com/boardgame/2719/connect-four)

## Specifications
Your AI must implement the [`fourinarowAI`](ai/driver/fourinarow_ai.go) interface. It will be passed a [`State`](ai/driver/state.go) struct, which contains a [`Board`](board.go) representing the current state of the game.

[`Board`](board.go) is a two-dimensional `int` slice with seven columns and six rows. Its values are 0-2:
```
0: empty
1: player 1
2: player 2
```

The starred square below is `[3][5]`:
```
-------------------------------
5 |   |   |   | * |   |   |   |
-------------------------------
4 |   |   |   |   |   |   |   |
-------------------------------
3 |   |   |   |   |   |   |   |
-------------------------------
2 |   |   |   |   |   |   |   |
-------------------------------
1 |   |   |   |   |   |   |   |
-------------------------------
0 |   |   |   |   |   |   |   |
-------------------------------
  | 0 | 1 | 2 | 3 | 4 | 5 | 6 |
```

## Moves
Your AI must return a [`Move`](move.go) containing the `Col` (0-6) of your next move. Note that row is unnecessary, as moves are inserted at the top of the board and drop through any empty space below.