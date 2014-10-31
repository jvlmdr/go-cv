package detect

type DetList interface {
	Len() int
	At(int) Det
}

type detListByScoreDesc struct{ DetList }

func (dets detListByScoreDesc) Less(i, j int) bool { return dets.At(i).Score > dets.At(j).Score }
func (dets detListByScoreDesc) Swap(i, j int)      { panic("read-only") }

// DetSlice wraps []Det to satisfy the DetList interface.
type DetSlice []Det

func (dets DetSlice) Len() int     { return len(dets) }
func (dets DetSlice) At(i int) Det { return dets[i] }
