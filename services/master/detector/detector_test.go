package detector

import (
    "github.com/Ygg-Drasill/DookieFilter/common/types"
    "github.com/Ygg-Drasill/DookieFilter/services/master/worker"
    "github.com/stretchr/testify/assert"
    "testing"
)

type MockWorker struct {
    *Worker
}

func getMockWorker() *MockWorker {
    return &MockWorker{
        Worker: &Worker{
            BaseWorker: worker.NewBaseWorker(
                nil,
                "detector"),
        },
    }
}

func (m *MockWorker) setupMockData() map[string]types.PlayerPosition {
    return map[string]types.PlayerPosition{
        "player1:10:0": {PlayerId: "player1", FrameIdx: 10, Position: types.Position{X: 3.0, Y: 4.0}},
        "player1:11:1": {PlayerId: "player1", FrameIdx: 11, Position: types.Position{X: 5.05, Y: 6.04}},
        "player2:10:0": {PlayerId: "player2", FrameIdx: 10, Position: types.Position{X: 5.0, Y: 6.0}},
        "player2:11:1": {PlayerId: "player2", FrameIdx: 11, Position: types.Position{X: 3.02, Y: 4.04}},
    }
}

func TestWorkerSwap(t *testing.T) {
    testCases := []struct {
        name       string
        p          map[string]types.PlayerPosition
        expectedXY map[string]types.Position
    }{
        {
            name: "Test Case 1",
            p: map[string]types.PlayerPosition{
                "player1:10:0": {PlayerId: "player1", FrameIdx: 10, Position: types.Position{X: 3.0, Y: 4.0}},
                "player1:11:1": {PlayerId: "player1", FrameIdx: 11, Position: types.Position{X: 5.05, Y: 6.04}},
                "player2:10:0": {PlayerId: "player2", FrameIdx: 10, Position: types.Position{X: 5.0, Y: 6.0}},
                "player2:11:1": {PlayerId: "player2", FrameIdx: 11, Position: types.Position{X: 3.02, Y: 4.04}},
            },
            expectedXY: map[string]types.Position{
                "player1:11:1": {X: 3.02, Y: 4.04},
                "player2:11:1": {X: 5.05, Y: 6.04},
            },
        },
        {
            name: "Test Case 2",
            p: map[string]types.PlayerPosition{
                "player1:10:0": {PlayerId: "player1", FrameIdx: 10, Position: types.Position{X: 3.0, Y: 4.0}},
                "player1:11:1": {PlayerId: "player1", FrameIdx: 11, Position: types.Position{X: 5.05, Y: 6.04}},
                "player2:10:0": {PlayerId: "player2", FrameIdx: 0, Position: types.Position{X: 0.0, Y: 0.0}},
                "player2:11:1": {PlayerId: "player2", FrameIdx: 11, Position: types.Position{X: 3.02, Y: 4.04}},
            },
            expectedXY: map[string]types.Position{
                "player1:11:1": {X: 3.02, Y: 4.04},
                "player2:11:1": {X: 5.05, Y: 6.04},
            },
        },
    }
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {

            mockWorker := getMockWorker()
            mockWorker.swap(true, tc.p)

            for key, expectedPos := range tc.expectedXY {
                if player, exists := tc.p[key]; exists {
                    assert.Equal(t, expectedPos, player.Position)
                } else {
                    t.Errorf("Player %s not found in map", key)
                }
            }
        })
    }
}

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
