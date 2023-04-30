# Reversi

## Overview
Players take turns placing discs on an 8x8 grid, flipping opponent discs bounded by each play.

## Reference Links
* [Wikipedia](https://en.wikipedia.org/wiki/Reversi)
* [BoardGameGeek](https://boardgamegeek.com/boardgame/2389/othello)

## House Rules
* The [Othello opening](https://en.wikipedia.org/wiki/Reversi#Rules) is used
* `Black` moves first

## Specifications
Your AI must implement the [`reversiAI`](ai/driver/reversi_ai.go) interface. It will be passed a [`State`](ai/driver/state.go) struct, which contains a [`Board`](board.go) representing the current state of the game.

[`Board`](board.go) is a two-dimensional `Disc` slice with eight columns and eight rows. Its values are `Black`, `White`, or `Empty`.

The starred square below is `[3][5]`:
```
-----------------------------------
7 |   |   |   |   |   |   |   |   |
-----------------------------------
6 |   |   |   |   |   |   |   |   |
-----------------------------------
5 |   |   |   | * |   |   |   |   |
-----------------------------------
4 |   |   |   |   |   |   |   |   |
-----------------------------------
3 |   |   |   |   |   |   |   |   |
-----------------------------------
2 |   |   |   |   |   |   |   |   |
-----------------------------------
1 |   |   |   |   |   |   |   |   |
-----------------------------------
0 |   |   |   |   |   |   |   |   |
-----------------------------------
  | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 |
```

## Moves
Your AI must return a [`Move`](move.go) containing the `[Col][Row]` of your next move.