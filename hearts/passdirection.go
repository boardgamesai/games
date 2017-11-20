package hearts

type PassDirection string

const (
	PassLeft   = PassDirection("left")
	PassAcross = PassDirection("across")
	PassRight  = PassDirection("right")
	PassNone   = PassDirection("none")
)

func (p PassDirection) Next() PassDirection {
	var direction PassDirection

	switch p {
	case PassLeft:
		direction = PassAcross
	case PassAcross:
		direction = PassRight
	case PassRight:
		direction = PassNone
	case PassNone:
		direction = PassLeft
	}

	return direction
}
