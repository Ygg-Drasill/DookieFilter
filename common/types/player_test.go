package types

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

var cases = map[string]SmallFrame{
	"random": {
		FrameIdx: 42,
		Players: []PlayerPosition{
			{
				Position: Position{
					X: rand.Float64(),
					Y: rand.Float64(),
				},
				FrameIdx:  42,
				PlayerNum: rand.Intn(100),
			},
		},
		Ball: Position{
			X: rand.Float64(),
			Y: rand.Float64(),
		},
	},
}

func TestSerializeFrame(t *testing.T) {
	for testCase, frame := range cases {
		t.Run(testCase, _testSerializeFrame(frame))
	}
}

func _testSerializeFrame(frame SmallFrame) func(t *testing.T) {
	return func(t *testing.T) {
		t.Parallel()
		var before = frame
		var after SmallFrame

		serialized := SerializeFrame(before)
		assert.NotEmptyf(t, serialized, "Expecetd non-empty serialization")
		after = DeserializeFrame(serialized)
		assert.Equal(t, before, after, "Expected frame to be equal before and after")
	}
}
