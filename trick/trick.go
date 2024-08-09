package trick

type Closed interface {
	Claims() []Claim // can be isolated from trick - claims work the same way in dopplekopf as in skat?? right??
	Winner() int
	Points() uint8
}
