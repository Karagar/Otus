package hw04_lru_cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("custom", func(t *testing.T) {
		l := NewList()

		l.PushBack("odin")    // [odin]
		first := l.Front()    // odin : check it already first
		l.Remove(first)       // [] check what if first and last was deleted
		l.PushBack("dva")     // [dva]
		l.PushFront("tri")    // [tri, dva]
		l.PushBack("chetyre") // [tri, dva, chetyre]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // dva
		require.Equal(t, "dva", middle.Value)
		l.Remove(middle) // [tri, chetyre]
		require.Equal(t, 2, l.Len())
		l.Remove(l.Back())
		l.MoveToFront(l.Back()) // [tri] check MoveToFront with 1 elem

		require.Equal(t, "tri", l.Front().Value)
		require.Equal(t, "tri", l.Back().Value)
	})
}
