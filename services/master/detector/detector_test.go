package detector

import (
    "github.com/Ygg-Drasill/DookieFilter/common/types"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestPositionProximity(t *testing.T) {
    // Test cases
    testCases := []struct {
        a, b     swapPlayer
        expected bool
    }{
        {
            a:        swapPlayer{player: types.PlayerPosition{Position: types.Position{X: 1.0, Y: 2.0}}},
            b:        swapPlayer{player: types.PlayerPosition{Position: types.Position{X: 1.0, Y: 2.0}}},
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
