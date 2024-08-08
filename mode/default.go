package mode

import (
	"fmt"

	"github.com/KonstantinGasser/playme/card"
)

type Card struct {
	kind  card.Kind
	value card.Value
}

func (c Card) Kind() card.Kind   { return c.kind }
func (c Card) Value() card.Value { return c.value }

func (c Card) IsTrump() bool {

	// Herz 10
	if c.Kind() == card.Herz && c.Value() == card.Ten {
		return true
	}

	// Any Unter || Ober
	switch c.Value() {
	case card.Unter, card.Ober:
		return true
	}

	// Any Schelle card is trump in a normal game
	if c.Kind() == card.Schellen {
		return true
	}

	return false
}

func (c Card) GreaterThan(other card.Mode) bool {

	// nothing beats an already played Herz 10
	if c.Kind() == card.Herz && c.Value() == card.Ten {
		return true // self >= other
	}

	if !c.IsTrump() && other.IsTrump() {
		return false // self < trump
	}

	if c.IsTrump() && !other.IsTrump() {
		return true // self (trump) >= non trump
	}

	if c.IsTrump() && other.IsTrump() {
		// unter && unter || ober && ober
		// => check for kind
		if (c.Value() == card.Unter && other.Value() == card.Unter) ||
			(c.Value() == card.Ober && other.Value() == card.Ober) {
			// assumes card.Kind ordered in correct order
			// TODO: maybe a test in the card package to ensure that?
			return c.Kind() >= other.Kind()
		}

		// unter && ober || ober && unter
		// => ober wins
		if (c.Value() == card.Unter && other.Value() == card.Ober) ||
			(c.Value() == card.Unter && other.Value() == card.Unter) {
			// assumes card.Value is orded in correct order. At least Ober > Unter
			// TODO: maybe a test in the card package to ensure that?
			return c.Value() >= other.Value()
		}

		// unter && schellen (non ober || unter)
		// => unter wins
		switch {
		case c.Value() == card.Unter, c.Value() == card.Ober: // self is unter || ober
			if other.Value() != card.Unter && other.Value() != card.Ober {
				// ok because checks before disallow that other card can
				// be Herz 10
				return true
			}
			// WARNING:
			// should be unreachable - if condition not meet will result
			// in unreachable panic at the end of the function

		}
		if (c.Value() == card.Unter || c.Value() == card.Ober) &&
			(other.Value() != card.Unter && other.Value() != card.Ober) {

		}

		// schellen && schellen
		// => check for value
	}

	// WARNING:
	// any combination of cards which are not machted in either of the above listed conditions
	// implies missing implemented rules and therefore no statement about self >= other can be made.
	// Hence, panic with info about the combination.
	panic(fmt.Sprintf("mode.default.GreaterThan: uncaught combination of cards! Self: %+v; Other: %+v", c, other))
}
