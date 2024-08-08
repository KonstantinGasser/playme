package mode

import "github.com/KonstantinGasser/playme/card"

type Card struct {
	kind  card.Kind
	value card.Value
}

func (c Card) Kind() card.Kind                  { return c.kind }
func (c Card) Value() card.Value                { return c.value }
func (c Card) Points() uint8                    { return 0 }
func (c Card) IsTrump() bool                    { return false }
func (c Card) GreaterThan(other card.Mode) bool { return false }
