package dice

import (
	"fmt"

	"github.com/boardgamesai/games/util"
)

type Dice[D any] struct {
	Values    []D // Current values of each die
	dieValues []D // Possible values for each die
}

// Convenience type for the common six-sided 1-6 dice
type Standard struct {
	*Dice[int]
}

func New[D any](numDice int, values []D) *Dice[D] {
	d := Dice[D]{
		Values:    make([]D, numDice),
		dieValues: values,
	}
	return &d
}

func NewNumeric(numDice, min, max int) *Dice[int] {
	vals := []int{}
	for i := min; i <= max; i++ {
		vals = append(vals, int(i))
	}

	return New(numDice, vals)
}

func NewStandard(numDice int) Standard {
	d := Standard{
		Dice: NewNumeric(numDice, 1, 6),
	}
	return d
}

func (d *Dice[D]) Roll() {
	for i := 0; i < len(d.Values); i++ {
		d.Values[i] = d.dieValues[util.RandInt(0, len(d.dieValues)-1)]
	}
}

func (d *Dice[D]) Add(value D) {
	d.Values = append(d.Values, value)
}

func (d *Dice[D]) Remove(index int) D {
	if index >= len(d.Values) || index < 0 {
		// This is programmer error, fail hard & fast
		panic(fmt.Sprintf("index %d out of range in Remove", index))
	}

	val := d.Values[index]
	d.Values = append(d.Values[:index], d.Values[index+1:]...)

	return val
}

func (d *Dice[D]) Count() int {
	if d == nil {
		return 0
	}
	return len(d.Values)
}

func (d *Dice[D]) String() string {
	return fmt.Sprintf("%s", d.Values)
}
