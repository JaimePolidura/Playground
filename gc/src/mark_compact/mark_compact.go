package markcompact

import "sync/atomic"

type MarkCompactGC struct {
	Data     [1048576]byte // 1MB
	NextFree atomic.Int32
}

func (this *MarkCompactGC) TryAllocate(size int32) int32 {
	if this.NextFree.Load() > 1048576 {
		return 0
	}

	newNextFree := this.NextFree.Add(int32(size))

	return newNextFree - size
}
