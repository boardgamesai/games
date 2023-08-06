# Hearts

## Overview
A trick-taking card game in the [Whist](https://en.wikipedia.org/wiki/Whist) family where you win by avoiding taking hearts.

## Reference Links
* [Wikipedia](https://en.wikipedia.org/wiki/Hearts_(card_game))
* [BoardGameGeek](https://boardgamegeek.com/boardgame/6887/hearts)

## House Rules
* Four players
* Jack of diamonds subtracts 10 points
* Match ends when any player reaches 100 points

## Specifications
Each match consists of a number of rounds. Each round consists of 13 tricks.

Your AI must implement the [`heartsAI`](ai/driver/hearts_ai.go) interface. It will be passed a [`State`](ai/driver/state.go) struct, which contains the following fields:

* `Position` - your position at the table, from 1-4
* `Trick` - the cards played thus far on the current trick
* `TrickCount` - the number of tricks played thus far in the current round, from 0-12
* `HeartsBroken` - indicates whether hearts have already been broken in this round
* `CurrentRound` - scores for the current round in progress
* `Scores` - scoring history per round, only gets written after a round is over
* `Hand` - your hand, cards are removed as they are played
* `PassDirection` - direction of passing in the current round, of type [`PassDirection`](passdirection.go)

## Moves
Your AI must return two different types of [moves](move.go):
1. `PassMove` - exactly three cards to pass from a newly-dealt hand, will not be called if pass direction is `PassNone`
1. `PlayMove` - a single card from your hand to play in the current trick