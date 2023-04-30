# Game of the Amazons

## Overview
On a 10x10 chessboard, players move four queens and shoot arrows to trap their opponents.

## Reference Links
* [Wikipedia](https://en.wikipedia.org/wiki/Game_of_the_Amazons)
* [BoardGameGeek](https://boardgamegeek.com/boardgame/2125/amazons)

## House Rules
* `White` moves first

## Specifications
Your AI must implement the [`amazonsAI`](ai/driver/amazons_ai.go) interface. It will be passed a [`State`](ai/driver/state.go) struct, which contains a [`Board`](board.go) representing the current state of the game.

[`Board`](board.go) is a two-dimensional `SpaceType` slice with ten columns and ten rows. Its values are `White`, `Black`, `Arrow`, or `Empty`.

The starred square below is `[3][5]`:
```
-------------------------------------------
9 |   |   |   |   |   |   |   |   |   |   |
-------------------------------------------
8 |   |   |   |   |   |   |   |   |   |   |
-------------------------------------------
7 |   |   |   |   |   |   |   |   |   |   |
-------------------------------------------
6 |   |   |   |   |   |   |   |   |   |   |
-------------------------------------------
5 |   |   |   | * |   |   |   |   |   |   |
-------------------------------------------
4 |   |   |   |   |   |   |   |   |   |   |
-------------------------------------------
3 |   |   |   |   |   |   |   |   |   |   |
-------------------------------------------
2 |   |   |   |   |   |   |   |   |   |   |
-------------------------------------------
1 |   |   |   |   |   |   |   |   |   |   |
-------------------------------------------
0 |   |   |   |   |   |   |   |   |   |   |
-------------------------------------------
  | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 |
```

## Moves
Your AI must return a [`Move`](move.go) containing three sets of coordinates:
* `From` - the queen to move
* `To` - the queen's destination
* `Arrow` - the arrow shot by the queen from its destination

These coordinates are all represented as a `Space`, which has `Col` and `Row` attributes corresponding to the [`Board`](board.go) above.