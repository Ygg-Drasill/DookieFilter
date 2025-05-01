package detector

import (
    "github.com/Ygg-Drasill/DookieFilter/common/types"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestPositionProximity(t *testing.T) {
    testCases := []struct {
        a, b     swapPlayer
        expected bool
    }{
        {
            a:        swapPlayer{player: types.PlayerPosition{Position: types.Position{X: 0.0, Y: 2.0}}},
            b:        swapPlayer{player: types.PlayerPosition{Position: types.Position{X: 0.05, Y: 2.05}}},
            expected: true},
        {
            a:        swapPlayer{player: types.PlayerPosition{Position: types.Position{X: 1.0, Y: 2.0}}},
            b:        swapPlayer{player: types.PlayerPosition{Position: types.Position{X: 10.0, Y: 2.0}}},
            expected: false},
    }

    for _, tc := range testCases {
        result := positionProximity(tc.a, tc.b)
        assert.Equal(t, tc.expected, result)
    }
}

func TestSwapPlayers(t *testing.T) {
    testCases := []struct {
        p                    map[string]types.PlayerPosition
        p1, p2               swapPlayer
        expectedA, expectedB types.Position
    }{
        {
            p: map[string]types.PlayerPosition{
                "player1": {Position: types.Position{X: 3.0, Y: 4.0}},
                "player2": {Position: types.Position{X: 1.0, Y: 2.0}}},
            p1: swapPlayer{
                key: "player1",
                player: types.PlayerPosition{
                    Position: types.Position{X: 3.0, Y: 4.0}}},
            p2: swapPlayer{
                key: "player2",
                player: types.PlayerPosition{
                    Position: types.Position{X: 1.0, Y: 2.0}}},
            expectedA: types.Position{X: 1.0, Y: 2.0},
            expectedB: types.Position{X: 3.0, Y: 4.0}},
    }
    for _, tc := range testCases {
        swapPlayers(tc.p, tc.p1, tc.p2)
        assert.Equal(t, tc.expectedA, tc.p[tc.p1.key].Position)
        assert.Equal(t, tc.expectedB, tc.p[tc.p2.key].Position)
    }
}
