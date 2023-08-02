# Liar's Dice

## Overview
A dice-rolling game where players bluff their way to be the last one standing.

## Reference Links
* [Wikipedia](https://en.wikipedia.org/wiki/Liar%27s_dice)
* [BoardGameGeek](https://boardgamegeek.com/boardgame/45/perudo)

## House Rules
* Four players
* Players start with five dice
* Follow the [Milton Bradley rules](../docs/liarsdice.pdf) (except for equal bid challenges, see below)
* On a challenge, if actual amount is:
    * Lower than the bid: bidder loses the difference
    * Higher than the bid: challenger loses the difference
    * Equal to the bid: challenger gives bidder one die
* Shown dice do not need to match the bid
* You cannot show all your dice - you must always have at least one die hidden

## Specifications
Your AI must implement the [`liarsdiceAI`](ai/driver/liarsdice_ai.go) interface. It will be passed a [`State`](ai/driver/state.go) struct, which contains the following fields:

* `Position` - your position at the table, from 1-4
* `Players` - all players (including you)
* `Dice` - your hidden dice
* `Bid` - the current bid, e.g. for `5 3s`, the bid is `3`
* `Quantity` - the current quantity, e.g. for `5 3s`, the quantity is `5`
* `Bidder` - the player who made the last bid
* `DiceCounts` - a map of how many dice each player has remaining
* `DiceShown` - a map of the dice each player is currently showing

## Moves
Your AI must return a [Move](move.go) in the following format:
* `Challenge` - set to `true` if you are challenging, no further fields need be set if so
* `Bid` - your new bid, e.g. for `5 3s`, the bid is `3`
* `Quantity` - your new quantity, e.g. for `5 3s`, the quantity is `5`
* `ShowDice` - any dice you wish to show, your remaining hidden dice will be re-rolled