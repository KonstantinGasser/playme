package main

import "fmt"

type Color int
type Value uint8

const (
	Eichel Color = iota
	Blatt
	Herz
	Shell

	Eight Value = iota
	Nine
	Ten
	Unter
	Ober
	King
	Ace
)

type Player struct {
	id   int
	hand [12]Card
	// a color present in this map indicates
	// that the player played a card other than the first one
	// (not counting for trump cards) and therefore no longer
	// has set color in his hand.
	// This is manily used to validate that the user confirms with
	// the playing rules within a tick
	missingColor map[Color]struct{}
}

func (player *Player) MarkAsMissing(c Color) {

	if _, ok := player.missingColor[c]; ok {
		return
	}

	player.missingColor[c] = struct{}{}
}

type Card struct {
	color    Color
	number   Value
	playedBy *Player
}

func (card Card) isTrump() bool { return false }

// type Mode = func(highCard Card, other Card) Card

type Mode interface {
	IsTrump(Card) bool
	BeatsHighCard(highcard Card, other Card) bool
	Points(Card) uint8
}

type Game struct {
	players               [4]Player
	firstTrickPlayerIndex int // on new game always player 0 == index = 0
	scores                uint8
	tricks                <-chan Trick
	events                []any
	mode                  Mode
}

func (game *Game) Run() {

	// when the game is started it knows the players and the player
	// that start the first round (_game.firstTrickPlayerIndex_).
	// Since we need to trace a chain of who won the last tick
	// we keep an initial reference and mutate whenever there
	// is a new winner after the evaluation
	lastTrickPlayerIndex := game.firstTrickPlayerIndex

	// start waiting for scanner input.
	// scanner sends 4 cards (one trick) at the time.
	for trick := range game.tricks {

		// who ever starts the trick will by default be the winner
		// of the trick until beaten by other cards/players
		currentWinnerIndex := lastTrickPlayerIndex

		for playerIndex, card := range trick.cards {

			// counting points is independent from who wins the
			// trick and can be counted directely without checks.
			game.scores += game.mode.Points(card)

			// check if new played cards rankes higher than current
			// high card and replace if so. if changes update current
			// winner of the trick
			if game.mode.BeatsHighCard(trick.highCard, card) {
				trick.highCard = card
				// index = 3 -> last card/player played winning card
				currentWinnerIndex = lastTrickPlayerIndex + playerIndex
			}

			// check if player violated trick contribution rules. If the player in a previous
			// trick claimed to no longer have a color and discarded (either trump or non trump)
			// card any further card played by the player __must__ not be listed in the player's
			// _missingColor_ map.
			if _, ok := game.players[lastTrickPlayerIndex+playerIndex].missingColor[card.color]; ok {
				// violation found by player.
			}

			// check if player _claims_ to no longer have trump cards or cards of a specific color.
			// This is relevant because the game needs to check for each played card if a player
			// forgot to correctly give a card.
			// The claim to no longer have a specific color is a beliefe which must hold true until
			// the end of the game without encountering any violations.
			// Check is irrelevant for the first played card of the trick.
			if playerIndex > 0 {
				switch {
				// 1. when the trick is kickedoff with _trump_ playes must serve trump.
				//    A player that fails to play trump _claims_ that it no longer has any trump cards.
				case game.mode.IsTrump(trick.openingCard) && !game.mode.IsTrump(card):
					fallthrough
				// 2. when the trick is kickedoff with __not__ a trump card but a player
				//    plays a trump card, again the user claims to no longer have the opening color.
				case !game.mode.IsTrump(trick.openingCard) && game.mode.IsTrump(card):
					fallthrough
				// 3. When neither the opening card __nor__ the played card is a trump card
				//    and the played card's color does not match the opening card's color,
				//    the user discarded one of its cards, again _claiming_ to no longer have the opening card's color
				case !game.mode.IsTrump(trick.openingCard) && !game.mode.IsTrump(card) && card.color != trick.openingCard.color:
					// player claims to no longer have either trump or a specific color
					game.players[lastTrickPlayerIndex+playerIndex].MarkAsMissing(trick.openingCard.color)
				}
			}
		}

		// after each trick update the index of the player
		// who won the previous trick / is starting the next trick.
		// By default winner of last trick will be initial winner of
		// next trick until beaten by another player
		lastTrickPlayerIndex = currentWinnerIndex
	}

}

type Trick struct {
	cards       [4]Card
	openingCard Card
	highCard    Card
}

func NewTrick(cards [4]Card) Trick {
	return Trick{
		cards:       cards,
		openingCard: cards[0],
		highCard:    Card{}, // should know the player
	}
}

func main() {

}
