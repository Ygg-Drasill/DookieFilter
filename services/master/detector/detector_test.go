package detector

import (
    "github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
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

func TestWorkerSwap(t *testing.T) {
    testCases := []struct {
        name       string
        p          map[string]types.PlayerPosition
        expectedXY map[string]types.Position
    }{
        {
            name: "Successful Swap",
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
            name: "Successful Swap with missing previous frame",
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
        {
            name: "Successful Swap with missing current frame",
            p: map[string]types.PlayerPosition{
                "player1:10:0": {PlayerId: "player1", FrameIdx: 10, Position: types.Position{X: 3.0, Y: 4.0}},
                "player1:11:1": {PlayerId: "player1", FrameIdx: 11, Position: types.Position{X: 5.05, Y: 6.04}},
                "player2:10:0": {PlayerId: "player2", FrameIdx: 10, Position: types.Position{X: 5.0, Y: 6.0}},
                "player2:11:1": {PlayerId: "player2", FrameIdx: 0, Position: types.Position{X: 0, Y: 0}},
            },
            expectedXY: map[string]types.Position{
                "player1:11:1": {X: 0.0, Y: 0.0},
                "player2:11:1": {X: 5.05, Y: 6.04},
            },
        },
        {
            name: "Unsuccessful Swap no match found",
            p: map[string]types.PlayerPosition{
                "player1:10:0": {PlayerId: "player1", FrameIdx: 10, Position: types.Position{X: 3.0, Y: 4.0}},
                "player1:11:1": {PlayerId: "player1", FrameIdx: 11, Position: types.Position{X: 6.05, Y: 6.04}},
                "player2:10:0": {PlayerId: "player2", FrameIdx: 10, Position: types.Position{X: 5.0, Y: 6.0}},
                "player2:11:1": {PlayerId: "player2", FrameIdx: 0, Position: types.Position{X: 8.0, Y: 4.0}},
            },
            expectedXY: map[string]types.Position{
                "player1:11:1": {X: 6.05, Y: 6.04},
                "player2:11:1": {X: 8.0, Y: 4.0},
            },
        },
        {
            name: "Successful 4 player Swap",
            p: map[string]types.PlayerPosition{
                "player1:10:0": {PlayerId: "player1", FrameIdx: 10, Position: types.Position{X: 3.0, Y: 4.0}},
                "player1:11:1": {PlayerId: "player1", FrameIdx: 11, Position: types.Position{X: 5.05, Y: 6.04}},
                "player2:10:0": {PlayerId: "player2", FrameIdx: 10, Position: types.Position{X: 2.0, Y: 7.0}},
                "player2:11:1": {PlayerId: "player2", FrameIdx: 11, Position: types.Position{X: 5.02, Y: -3.04}},
                "player3:10:0": {PlayerId: "player3", FrameIdx: 10, Position: types.Position{X: 5.0, Y: -3.0}},
                "player3:11:1": {PlayerId: "player3", FrameIdx: 11, Position: types.Position{X: 2.05, Y: 6.94}},
                "player4:10:0": {PlayerId: "player4", FrameIdx: 10, Position: types.Position{X: 5.0, Y: 6.0}},
                "player4:11:1": {PlayerId: "player4", FrameIdx: 0, Position: types.Position{X: 3.02, Y: 4.04}},
            },
            expectedXY: map[string]types.Position{
                "player1:11:1": {X: 3.02, Y: 4.04},
                "player2:11:1": {X: 2.05, Y: 6.94},
                "player3:11:1": {X: 5.02, Y: -3.04},
                "player4:11:1": {X: 5.05, Y: 6.04},
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

func TestJump(t *testing.T) {
    testCases := []struct {
        name        string
        p           map[string]types.PlayerPosition
        stateBuffer []*types.SmallFrame
        expected    bool
    }{
        {
            name: "Player Duplicate Coords Detected",
            p: map[string]types.PlayerPosition{
                "player1:10:0": {PlayerId: "player1", FrameIdx: 10, Position: types.Position{X: 3.0, Y: 4.0}},
                "player1:11:1": {PlayerId: "player1", FrameIdx: 11, Position: types.Position{X: 5.0, Y: 6.0}},
            },
            stateBuffer: []*types.SmallFrame{
                {FrameIdx: 10, Players: []types.PlayerPosition{
                    {PlayerId: "player2", FrameIdx: 10, Position: types.Position{X: 5.0, Y: 6.0}},
                }},
            },
            expected: true},
        {
            name: "Player Duplicate Coords Not Detected",
            p: map[string]types.PlayerPosition{
                "player1:10:0": {PlayerId: "player1", FrameIdx: 10, Position: types.Position{X: 3.0, Y: 4.0}},
                "player1:11:1": {PlayerId: "player1", FrameIdx: 11, Position: types.Position{X: 5.0, Y: 6.0}},
            },
            stateBuffer: []*types.SmallFrame{
                {FrameIdx: 10, Players: []types.PlayerPosition{
                    {PlayerId: "player2", FrameIdx: 10, Position: types.Position{X: 8.0, Y: 6.0}},
                }},
            },
            expected: false,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            buffer := pringleBuffer.New[types.SmallFrame](10)
            for _, frame := range tc.stateBuffer {
                buffer.Insert(*frame)
            }
            w := getMockWorker()
            w.stateBuffer = buffer
            r := w.jump(tc.p)
            assert.Equal(t, tc.expected, r)
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
