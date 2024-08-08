package card

type Kind int
type Value uint8

const (
	Eichel Kind = iota
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

type Mode interface {
	Kind() Kind
	Value() Value
	Points() uint8
	IsTrump() bool
	GreaterThan(other Mode) bool
}
