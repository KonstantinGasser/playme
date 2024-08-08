package game

import (
	"github.com/KonstantinGasser/playme/card"
	"github.com/KonstantinGasser/playme/trick"
)

type Scanner interface {
	Scan() [4]card.Mode
}

type Dopplekopf struct {
	// a call to scan should return the next 4 playing cards
	// to be evaluated together.
	scanner Scanner
	// maintains state about each player including claims of a player.
	players [4]Player
	// indicates the first player of the first trick played.
	// This is used to reconstruct a causal
	firstTrickPlayerIndex int // on new game always player 0 == index = 0
	tricks                <-chan trick.Open
	events                []any
}

func NewDoppleKopf(scanner Scanner, players [4]Player) *Dopplekopf {
	return &Dopplekopf{
		scanner:               scanner,
		players:               players,
		firstTrickPlayerIndex: 0,
		tricks:                make(<-chan trick.Open),
		events:                []any{},
	}
}

// Listen is a background process which allows
// externals to send card information via a defined
// interface (network, pipes, etc.) to the game engine.
// Cards will be processed and evaluated in the order received.
// Thus it is the responsability of the external to ensure no
// rule violations introduced by re-ordering occure.
func (game Dopplekopf) Listen() {}

// Process is a blocking call that returns either on
// unrecoverable errors or once a game is determined
// to be finished.
//
// A call to the function will start the generic processing
// and evaluation of received tricks. Which mode is played
// is not important and abstracted by what ever mode is reflected
// by the implemented cards.
func (game *Dopplekopf) Process() {

	// when the game is started it knows the players and the player
	// that start the first round (_game.firstTrickPlayerIndex_).
	// Since we need to trace a chain of who won the last trick
	// we keep an initial reference and mutate whenever there
	// is a new winner after the evaluation.
	// Knowing the last winner implies knowing who will start
	// the next trick. For the the trick that is since evaulations
	// such as claims or party matching is tightly coupled to
	// who played the card.
	previousTrickWinnerIndex := game.firstTrickPlayerIndex

	// start waiting for scanner input.
	// scanner sends 4 cards (one trick) at the time.
	for trick := range game.tricks {

		// who ever starts the trick will by default be the winner
		// of the trick until beaten by other cards/players
		// trick.openingPlayerIndex = lastTrickPlayerIndex
		closedTrick := trick.Replay(previousTrickWinnerIndex)

		// check if player violated trick contribution rules. If the player in a previous
		// trick claimed to no longer have a color and discarded (either trump or non trump)
		// card any further card played by the player __must__ not be listed in the player's
		// _missingColor_ map.
		game.FindViolation(closedTrick)

		// after each trick update the index of the player
		// who won the previous trick / is starting the next trick.
		// By default winner of last trick will be initial winner of
		// next trick until beaten by another player
		previousTrickWinnerIndex = closedTrick.WinnerIndex()
	}
}

func (game *Dopplekopf) FindViolation(closed trick.Closed) {
	for _, claim := range closed.Claims() {
		if _, ok := game.players[claim.By()].missingKind[claim.Kind()]; ok {
			// violation found by player.
			// TODO: What now? What should we do?
		}
	}
}
