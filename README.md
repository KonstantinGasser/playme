
#### Preconditions
- for each round the player order _must_ remain the same
- the machine needs to know which player started


#### Single Trick evaluations

__Example__
Player 0 starts the trick. Playing order read from left to right, convinentily with increasing player IDs (only in these examples)

_Single colour tick_
Trick cards: 0: Ace Green, 1: 9 Green, 2: 10 Green, 3: King Green
Winner: Player 0

_Single over rule_
Tick cards: 0: Ace Green, 1: 9 Green, 2: 10 Green, 3: Ace Trump
Winner: Player 3

_Multiple over rules_
Trick cards: 0: Ace Heart, 1: Unter Green, 2: 9 Heart, 3: Ober Eichel
Winner: _after scanning player 0 and player 1_: player 1 beats player 0. _after full scan player 3 beats player 1:_ player 3 winns
_Note:_ Player 1 and player 3 do not have green anymore which needs to be noted to later check if a player violates the game rules


_Multiple over rules_
Trick cards: 0: Ace Heart, 1: Unter Green, 2: 9 Heart, 3: Unter Green
Winner: Player 1 winns, player 3's card not higher

#### Code Structure

__Round__
- maintains state 
    - player
        - figure out which players are Re party
    - scores
    - non_present_colors_in_player_hand
    - events (fox, dopplekopf, karlchen macht den letzten)
 - offers pub/sub to Display interface

__Tick__
- maintains state about single tick
    - current high card and played by
- checks for single tick validitiy / game violations
- communicates events happened in a tick (streams to round state)
- communicates scores for player (streams to round state)

__Interface (HTTP or else)__
- provide data for display

Test change for submodule
