package trick

import "github.com/KonstantinGasser/playme/card"

type Claim struct {
	cardKind    card.Kind
	playerIndex int
}

func (claim Claim) By() int {
	return claim.playerIndex
}

func (claim Claim) Kind() card.Kind {
	return claim.cardKind
}
