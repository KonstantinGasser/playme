package game

import "github.com/KonstantinGasser/playme/card"

type Party int

const (
	UnknownParty Party = iota
	Re
	Contra
)

type Player struct {
	id   int
	hand [12]card.Card
	// lazily evaluated throughout the round.
	// __Should__ not be "Unknown" after the round
	// is finished.
	// A player in Re implies that used __must__ have had
	// the Card "Eichel Ober" in a normal game mode else
	// dependend on the game mode.
	party Party
	// a color present in this map indicates
	// that the player played a card other than the first one
	// (not counting for trump cards) and therefore no longer
	// has set color in his hand.
	// This is manily used to validate that the user confirms with
	// the playing rules within a tick
	missingKind map[card.Kind]struct{}
}
