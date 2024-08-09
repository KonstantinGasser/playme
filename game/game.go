package game

import (
	"github.com/KonstantinGasser/playme/card"
	"github.com/KonstantinGasser/playme/trick"
)

// TODO: Not sure about that one..
// Will be figured out once we start with
// the implementation of the scanner -> engine
type Scanner interface {
	Scan() [4]card.Card
}

type Replayer interface {
	Replay(previousWinnerIndex int) trick.Played
}

// something that can be shown on the display
// when an violation is recoded
type ErrRuleViolation struct{}

type Recorder interface {
	Apply(trick.Claim) *ErrRuleViolation
	Add(uint8) error
}

type Game[T Replayer, R Recorder] struct {
	// a call to scan should return the next 4 playing cards
	// to be evaluated together.
	scanner  Scanner
	recorder Recorder
	// indicates the first player of the first trick played.
	// This is used to reconstruct a causal
	firstTrickPlayerIndex int // on new game always player 0 == index = 0
	tricks                <-chan T
	events                []any
}

func NewDoppleKopf[T Replayer, R Recorder](scanner Scanner, recorder Recorder) *Game[T, R] {
	return &Game[T, R]{
		scanner:               scanner,
		recorder:              recorder,
		firstTrickPlayerIndex: 0,
		tricks:                make(<-chan T),
		events:                []any{},
	}
}

// Listen is a background process which allows
// externals to send card information via a defined
// interface (network, pipes, etc.) to the game engine.
// Cards will be processed and evaluated in the order received.
// Thus it is the responsability of the external to ensure no
// rule violations introduced by re-ordering occure.
func (game Game[T, R]) Listen() {}

// Process is a blocking call that returns either on
// unrecoverable errors or once a game is determined
// to be finished.
//
// A call to the function will start the generic processing
// and evaluation of received tricks. Which mode is played
// is not important and abstracted by what ever mode is reflected
// by the implemented cards.
func (game *Game[T, R]) Process() {

	// when the game is started it knows the players and the player
	// that start the first round (_game.firstTrickPlayerIndex_).
	// Since we need to trace a chain of who won the last trick
	// we keep an initial reference and mutate whenever there
	// is a new winner after the evaluation.
	// Knowing the last winner implies knowing who will start
	// the next trick. For the the trick that is since evaulations
	// such as claims or party matching is tightly coupled to
	// who played the card.
	// TODO: maybe this should also be moved into the recorder?
	// In the recorder this information might be required anyway whereas
	// here it is not really except for the Replayer. But then again
	// maybe Recorder should be generic over a Replayer?
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

		// forward counted points in the trick to recorder.
		game.recorder.Add(closedTrick.Points())

		// after each trick update the index of the player
		// who won the previous trick / is starting the next trick.
		// By default winner of last trick will be initial winner of
		// next trick until beaten by another player
		previousTrickWinnerIndex = closedTrick.Winner()
	}
}

func (game *Game[T, R]) FindViolation(ct trick.Played) {
	for _, claim := range ct.Claims() {
		if err := game.recorder.Apply(claim); err != nil {
			// TODO: figure out what to do when violation is found
		}
	}
}
