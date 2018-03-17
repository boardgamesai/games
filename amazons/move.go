package amazons

import "fmt"

type Move struct {
	From  Space
	To    Space
	Arrow Space
}

func (m Move) String() string {
	return fmt.Sprintf("from %s to %s arrow %s", m.From, m.To, m.Arrow)
}
