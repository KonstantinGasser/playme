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

func (c Card) GreaterThan(other card.Card) bool {

	// nothing beats an already played Herz 10
	if c.Kind() == card.Herz && c.Value() == card.Ten {
		return true // self >= other
	}

	// inverse: input is Herz 10 and I am not
	if c.Kind() != card.Herz && c.Value() != card.Ten && other.Kind() == card.Herz && other.Value() == card.Ten {
		return false
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
			// TODO: maybe a test in the `card` package to ensure that?
			return c.Kind() >= other.Kind()
		}

		// unter && ober || ober && unter
		// => ober wins
		if (c.Value() == card.Unter && other.Value() == card.Ober) ||
			(c.Value() == card.Unter && other.Value() == card.Unter) {
			// assumes card.Value is orded in correct order. At least Ober > Unter
			// TODO: maybe a test in the `card` package to ensure that?
			return c.Value() >= other.Value()
		}

		// unter && schellen (non ober || unter)
		// => unter wins
		if c.Value() == card.Unter || c.Value() == card.Ober { // self is unter || ober
			if other.Value() != card.Unter && other.Value() != card.Ober {
				// ok because checks before disallow that other card can
				// be Herz 10
				return true
			}

			// WARNING:
			// any other combination is an undeclared one which is undefined behaviour
			goto unreachable

		}
		if (c.Value() == card.Unter || c.Value() == card.Ober) &&
			(other.Value() != card.Unter && other.Value() != card.Ober) {
			// TODO: I guess we can return true because if that is the case
			// the other card must be a lower trump or a non trump card?
			// Herz 10 check already happend before, right
			return true

		}

		// schellen && schellen
		// => check for value
		if c.Kind() == card.Schellen && other.Kind() == card.Schellen {
			return c.Value() >= other.Value()
		}

		// WARNING:
		// any other combination is an undeclared one which is undefined behaviour
		goto unreachable
	}

	// both cards are non trump cards if they don't have the same kind
	// it is assumed that the compared to card (self) is the dominant card
	// and is >= other.
	// otherwise same kind comparsion by value.
	if !c.IsTrump() && !other.IsTrump() {
		if c.Kind() != other.Kind() {
			return true
		}

		// NOTE: assumes card.Values are orded in correct order
		return c.Value() >= other.Value()
	}

	// WARNING:
	// any combination of cards which are not machted in either of the above listed conditions
	// implies missing implemented rules and therefore no statement about self >= other can be made.
	// Hence, panic with info about the combination.
unreachable:
	panic(fmt.Sprintf("mode.default.GreaterThan: uncaught combination of cards! Self: %+v; Other: %+v", c, other))
}
