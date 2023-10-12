package teldrive

func (d *Teldrive) int64min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
