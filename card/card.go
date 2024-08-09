package card

type Kind int
type Value uint8

const (
	Eichel Kind = iota
	Blatt
	Herz
	Schellen

	Eight Value = 8
	Nine  Value = 9
	Ten   Value = 10
	Unter Value = 2
	Ober  Value = 3
	King  Value = 4
	Ace   Value = 11
)

type Card interface {
	Kind() Kind
	Value() Value
	IsTrump() bool
	GreaterThan(other Card) bool
	// IndicatesAffiliation() bool // not sure thou??
}
