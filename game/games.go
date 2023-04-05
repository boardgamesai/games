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

type Links map[LinkType]string
type LinkType string

var (
	Wikipedia     = LinkType("wikipedia")
	BoardGameGeek = LinkType("boardgamegeek")
)

type MetaData struct {
	Name        Name
	DisplayName string
	NumPlayers  int
	HasScore    bool
	HasTies     bool
	Description string
	Links       Links
}

var Data = map[Name]MetaData{
	Amazons: {
		Name:        Amazons,
		NumPlayers:  2,
		HasScore:    false,
		HasTies:     false,
		DisplayName: "Game of the Amazons",
		Description: "On a 10x10 chessboard, players move four queens and shoot arrows to trap their opponents.",
		Links: Links{
			Wikipedia:     "https://en.wikipedia.org/wiki/Game_of_the_Amazons",
			BoardGameGeek: "https://boardgamegeek.com/boardgame/2125/amazons",
		},
	},
	FourInARow: {
		Name:        FourInARow,
		DisplayName: "Four-in-a-Row",
		NumPlayers:  2,
		HasScore:    false,
		HasTies:     true,
		Description: "Players take turns dropping discs into a 6x7 grid, trying to get four in a row.",
		Links: Links{
			Wikipedia:     "https://en.wikipedia.org/wiki/Connect_Four",
			BoardGameGeek: "https://boardgamegeek.com/boardgame/2719/connect-four",
		},
	},
	Hearts: {
		Name:        Hearts,
		DisplayName: "Hearts",
		NumPlayers:  4,
		HasScore:    true,
		HasTies:     true,
		Description: "A trick-taking card game in the Whist family where you win by avoiding taking hearts.",
		Links: Links{
			Wikipedia:     "https://en.wikipedia.org/wiki/Hearts_(card_game)",
			BoardGameGeek: "https://boardgamegeek.com/boardgame/6887/hearts",
		},
	},
	Reversi: {
		Name:        Reversi,
		DisplayName: "Reversi",
		NumPlayers:  2,
		HasScore:    true,
		HasTies:     true,
		Description: "Players take turns placing discs on an 8x8 grid, flipping opponent discs bounded by each play.",
		Links: Links{
			Wikipedia:     "https://en.wikipedia.org/wiki/Reversi",
			BoardGameGeek: "https://boardgamegeek.com/boardgame/2389/othello",
		},
	},
	TicTacToe: {
		Name:        TicTacToe,
		DisplayName: "Tic-Tac-Toe",
		NumPlayers:  2,
		HasScore:    false,
		HasTies:     true,
		Description: "X and O take turns on a 3x3 grid, trying to get three in a row.",
		Links: Links{
			Wikipedia:     "https://en.wikipedia.org/wiki/Tic-tac-toe",
			BoardGameGeek: "https://boardgamegeek.com/boardgame/11901/tic-tac-toe",
		},
	},
	UltTicTacToe: {
		Name:        UltTicTacToe,
		DisplayName: "Ultimate Tic-Tac-Toe",
		NumPlayers:  2,
		HasScore:    false,
		HasTies:     true,
		Description: "A variation on tic-tac-toe where the goal is to win a grid comprised of smaller subgrids.",
		Links: Links{
			Wikipedia:     "https://en.wikipedia.org/wiki/Ultimate_tic-tac-toe",
			BoardGameGeek: "https://boardgamegeek.com/boardgame/9898/tic-tac-toe-times-10",
		},
	},
}
