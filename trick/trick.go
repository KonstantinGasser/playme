package trick

import "github.com/KonstantinGasser/playme/card"

type Open struct {
	cards              [4]card.Mode
	openingCard        card.Mode
	openingPlayerIndex int
	currentWinnerIndex int
	// claims express that a player claims a card
	// or in general did not serve the opening color and thus
	// no longer has the opening color. This is expressed as
	// belief since it cannot tell if that is a violation or not.
	// The collected claims can by used to evaluate violations
	// outside of the Trick
	claims   [3]Claim
	highCard card.Mode
	points   uint8
}

func NewTrick(openingPlayerIndex int, cards [4]card.Mode) Open {
	return Open{
		cards:              cards,
		openingCard:        cards[0],
		openingPlayerIndex: openingPlayerIndex,
		currentWinnerIndex: openingPlayerIndex,
		claims:             [3]Claim{},
		highCard:           nil, // should know the player
		points:             0,
	}
}

func (open Open) Close() Closed {
	return Closed{
		cards:             open.cards,
		playerWinnerIndex: open.currentWinnerIndex,
		claims:            open.claims,
		points:            open.points,
		events:            []any{},
	}
}

func (trick *Open) Replay(previousWinnerIndex int) Closed {

	trick.openingPlayerIndex = previousWinnerIndex

	for playerIndex, card := range trick.cards {

		// counting points is independent from who wins the
		// trick and can be counted directely without checks.
		trick.points += uint8(card.Value())

		// check if new played cards rankes higher than current
		// high card and replace if so. if changes update current
		// winner of the trick
		if card.GreaterThan(trick.highCard) { // mode.BeatsHighCard(trick.highCard, card) {
			trick.highCard = card
			// index = 3 -> last card/player played winning card
			trick.currentWinnerIndex = trick.openingPlayerIndex + playerIndex
		}

		// in normal game mode the card "Eichel Ober" played signales that the player
		// is part of the Re party. This, however, is game dependend and might be pre-defined
		// when selecting the game mode. Still we need to evaluate if a played card indicates
		// an affiliation to the Re party.
		// NOTE: players of the Contra party are passively determined as only the absents of a
		// card indicates their party.
		//
		// TODO: code goes her

		// check if player _claims_ to no longer have trump cards or cards of a specific color.
		// This is relevant because the game needs to check for each played card if a player
		// forgot to correctly give a card.
		// The claim to no longer have a specific color is a beliefe which must hold true until
		// the end of the game without encountering any violations.
		// Check is irrelevant for the first played card of the trick.
		if playerIndex > 0 {
			switch {
			// 1. when the trick is kickedoff with _trump_, playes must serve trump.
			//    A player that fails to play trump _claims_ that it no longer has any trump cards.
			case trick.openingCard.IsTrump() && !card.IsTrump():
				fallthrough
			// 2. when the trick is kickedoff with __not__ a trump card but a player
			//    plays a trump card, again the user claims to no longer have the opening color.
			case !trick.openingCard.IsTrump() && card.IsTrump():
				fallthrough
			// 3. When neither the opening card __nor__ the played card is a trump card
			//    and the played card's color does not match the opening card's color,
			//    the user discarded one of its cards, again _claiming_ to no longer have the opening card's color
			case !trick.openingCard.IsTrump() && !card.IsTrump() && card.Kind() != trick.openingCard.Kind():
				// player claims to no longer have either trump or a specific color.
				//
				// NOTE: index in the beliefs array is independent from player indecies.
				// it is only used to store beliefs which can be iterated over on game level
				trick.claims[playerIndex-1] = Claim{
					// color the player claims to no longer have
					cardKind: trick.openingCard.Kind(),
					// player index claming
					playerIndex: trick.openingPlayerIndex + playerIndex,
				}
			}
		}
	}

	return trick.Close()
}

type Closed struct {
	cards             [4]card.Mode
	playerWinnerIndex int
	claims            [3]Claim
	points            uint8
	events            []any
}

func (closed Closed) Claims() [3]Claim {
	return closed.claims
}

func (closed Closed) WinnerIndex() int {
	return closed.playerWinnerIndex
}
