package game

type Name string

const (
	Amazons      = Name("amazons")
	FourInARow   = Name("fourinarow")
	Hearts       = Name("hearts")
	Reversi      = Name("reversi")
	TicTacToe    = Name("tictactoe")
	UltTicTacToe = Name("ulttictactoe")
)

type MetaDataEntry struct {
	NumPlayers int
	HasScore   bool
}

var MetaData = map[Name]MetaDataEntry{
	Amazons: {
		NumPlayers: 2,
		HasScore:   false,
	},
	FourInARow: {
		NumPlayers: 2,
		HasScore:   false,
	},
	Hearts: {
		NumPlayers: 4,
		HasScore:   true,
	},
	Reversi: {
		NumPlayers: 2,
		HasScore:   true,
	},
	TicTacToe: {
		NumPlayers: 2,
		HasScore:   false,
	},
	UltTicTacToe: {
		NumPlayers: 2,
		HasScore:   false,
	},
}
