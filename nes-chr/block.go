package neschr

type BlockMetadata struct {
	Size     int
	Position int
}

func (bm *BlockMetadata) EndPosition() int {
	return bm.Position + bm.Size
}
