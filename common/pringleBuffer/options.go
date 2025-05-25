package pringleBuffer

func WithOnPopTail[TElement PringleIndexable](fun func(TElement)) func(pb *PringleBuffer[TElement]) {
	return func(pb *PringleBuffer[TElement]) {
		pb.onTailPop = fun
	}
}
