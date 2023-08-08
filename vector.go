package vecotor

import (
	"math"
	"sync"
	"unsafe"
)

const (
	initBlockLenBytes = 256
	regroupThreshold  = 0.2 // 20%
)

type Vector[T any] struct {
	a                [][]T
	blockLen         int
	curBlock, curPos int
	regroupCountDown int

	len int
	mx  sync.Mutex
}

func New[T any](a ...T) *Vector[T] {
	v := &Vector[T]{
		curBlock: 0,
		curPos:   0,
		len:      0,
	}

	v.Add(a...)
	return v
}

func (v *Vector[T]) Add(a ...T) {
	if len(a) == 0 {
		return
	}

	v.regroupCountDown--
	if v.regroupCountDown <= 0 {
		v.curBlock, v.curPos, v.blockLen, v.regroupCountDown = v.regroupBlocks(len(a))
	}

	v.addElems(a)
}

func (v *Vector[T]) Get(n int) (T, bool) {
	var (
		illegalIdx        = n < 0
		idxOutOfTheBorder = v.len == 0 || n > v.len
	)
	if illegalIdx || idxOutOfTheBorder {
		var t T
		return t, false
	}

	i, j := n/v.blockLen, n%v.blockLen
	return v.a[i][j], true
}

func (v *Vector[T]) Slice() []T {
	s := make([]T, v.len)
	si := 0
	for i := range v.a {
		copy(s[si:], v.a[i])
		si += len(v.a[i])
	}
	return s
}

func (v *Vector[T]) PushFront(a ...T) {
}

func (v *Vector[T]) Len() int {
	return v.len
}

func (v *Vector[T]) Find(a T) bool {
	return false
}

func (v *Vector[T]) FindMany(a ...T) []bool {
	return nil
}

func (v *Vector[T]) regroupBlocks(length int) (cb int, pos int, blockLen int, countDown int) {
	expectedBlkLen := v.expectedBlockLen(length)

	needRegroup := float64(v.blockLen)*(1.0+regroupThreshold) < float64(expectedBlkLen)
	if needRegroup {
		newTotalLen := length + v.len
		return v.regroup(newTotalLen, expectedBlkLen)
	}

	return v.curBlock, v.curPos, v.blockLen, v.regroupCountDown
}

func (v *Vector[T]) expectedBlockLen(length int) int {
	ebl := math.Ceil(math.Sqrt(float64(length + v.len)))

	var t T
	minBlockLen := float64(initBlockLenBytes) / float64(unsafe.Sizeof(t))

	if ebl < minBlockLen {
		ebl = minBlockLen
	}
	return int(ebl)
}

func (v *Vector[T]) regroup(newTotalLen, expectedBlkLen int) (cb int, pos int, blockLen int, countDown int) {
	upperSliceLen := int(math.Ceil(float64(newTotalLen) / float64(expectedBlkLen)))
	newA := make([][]T, 0, upperSliceLen)
	ni, nj := 0, 0

	// TODO: copy should be faster on big data amounts but it's hard to implement
	for i := range v.a {
		for j := range v.a[i] {
			for ni >= len(newA) {
				newA = append(newA, make([]T, expectedBlkLen))
			}
			newA[ni][nj] = v.a[i][j]
			nj++
			if nj == expectedBlkLen {
				ni++
				nj = 0
			}
		}
	}

	if ni < len(newA) {
		newA[ni] = newA[ni][:nj]
	}
	v.a = newA

	if nj+1 >= expectedBlkLen {
		return ni + 1, 0, expectedBlkLen, expectedBlkLen * 3
	}
	return ni, nj + 1, expectedBlkLen, expectedBlkLen * 3
}

func (v *Vector[T]) addElems(a []T) {
	for i := range a {
		for v.curBlock >= len(v.a) {
			v.a = append(v.a, make([]T, 0, v.blockLen))
		}

		v.a[v.curBlock] = append(v.a[v.curBlock], a[i])

		v.curPos++
		if v.curPos == v.blockLen {
			v.curPos = 0
			v.curBlock++
		}
	}
	v.regroupCountDown -= len(a)
	v.len += len(a)
}
