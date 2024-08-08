package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrickReplay(t *testing.T) {

	tt := []struct {
		name            string
		trick           OpenTrick
		wantWinnerIndex int
		wantDiscarded   [3]Claims
		points          uint8
	}{
		{
			name: "All same color",
			trick: NewTrick(0, [4]Card{
				{
					color:  Herz,
					number: Nine,
				},
				{
					color:  Herz,
					number: Eight,
				},
				{
					color:  Herz,
					number: King,
				},
				{
					color:  Herz,
					number: Ace,
				},
			}),
			wantWinnerIndex: 3,           // ace beats all
			wantDiscarded:   [3]Claims{}, // should be empty
			points:          15,
		},
	}

	for _, tc := range tt {

		closedTrick := tc.trick.Replay(0, nil)

		if !assert.Equal(t, tc.points, closedTrick.points) {
			t.Fatalf("%s: wanted: %d points; got: %d points", tc.name, tc.points, closedTrick.points)
		}
	}
}
