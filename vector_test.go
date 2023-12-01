package vecotor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Vector(t *testing.T) {
	t.Parallel()

	t.Run("happy", func(t *testing.T) {
		t.Parallel()

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
		t.Parallel()

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

func Test_VectorAdd(t *testing.T) {
	t.Parallel()

	t.Run("add empty", func(t *testing.T) {
		t.Parallel()

		vec := New(1, 2, 3)
		vec.Add()
		require.Equal(t, vec.Slice(), []int{1, 2, 3})
	})
}

func Test_VectorContains(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		vec := New(1, 2, 3)
		require.True(t, vec.Contains(Eq[int], 1))
		require.True(t, vec.Contains(Eq[int], 2))
		require.True(t, vec.Contains(Eq[int], 3))
	})
	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		vec := New(1, 2, 3)
		require.False(t, vec.Contains(Eq[int], 0))
		require.False(t, vec.Contains(Eq[int], 4))
		require.False(t, vec.Contains(Eq[int], -3))
	})
}

func Test_VectorContainsMany(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		vec := New(1, 2, 3)
		require.Equal(t, vec.ContainsMany(Eq[int], 1, 2, 3), []bool{true, true, true})
	})
	t.Run("not_found", func(t *testing.T) {
		t.Parallel()

		vec := New(1, 2, 3)
		require.Equal(t, vec.ContainsMany(Eq[int], 0, 4, -1), []bool{false, false, false})
	})
}
