package vecotor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Vector(t *testing.T) {
	t.Parallel()

	t.Run("happy", func(t *testing.T) {
		vec := New([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}...)
		vec.Add(222)
		vec.Add(223)
		require.Equal(t, 20, vec.Len())
		val, found := vec.Get(0)
		require.True(t, found)
		require.Equal(t, 1, val)
		val, found = vec.Get(-1)
		require.False(t, found)
		require.Equal(t, 0, val)
		val, found = vec.Get(19)
		require.True(t, found)
		require.Equal(t, 223, val)
		s := vec.Slice()
		require.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 222, 223}, s)
		require.Equal(t, 20, cap(s))
	})

	t.Run("large", func(t *testing.T) {
		vec := New([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}...)
		for i := 0; i < 10000; i++ {
			vec.Add([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}...)
		}
		require.Equal(t, 180018, vec.len)
		require.Equal(t, 180018, len(vec.Slice()))
		val, ok := vec.Get(18000)
		require.True(t, ok)
		require.Equal(t, 1, val)
		require.Equal(t, 1, vec.Slice()[18000])
	})
}
