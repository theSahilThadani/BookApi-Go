package store

import (
	"bookapi/internal/model"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
)

func mkCreateSharded(i int) model.CreateBookRequest {
	return model.CreateBookRequest{
		Title:       "Title-" + strconv.Itoa(i),
		Author:      "Author",
		ISBN:        fmt.Sprintf("ISBN-%d", i),
		PublishYear: 2024,
	}
}

func BenchmarkSharded_CreateParallel(b *testing.B) {
	s := NewShardedBookStore(32)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			s.Create(mkCreateSharded(i))
			i++
		}
	})
}

func BenchmarkSharded_MixedParallel(b *testing.B) {
	s := NewShardedBookStore(32)
	var idCounter uint64
	// pre-populate
	for i := 0; i < 1000; i++ {
		s.Create(mkCreateSharded(i))
		atomic.AddUint64(&idCounter, 1)
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			n := int(atomic.AddUint64(&idCounter, 1) % 4)
			switch n {
			case 0:
				s.Create(mkCreateSharded(int(idCounter)))
			case 1:
				s.GetById(fmt.Sprintf("BOOK-%d", (idCounter%1000)+1))
			case 2:
				var str = "Update"
				req := model.UpdateBookRequest{Title: &str}
				s.Update(fmt.Sprintf("BOOK-%d", (idCounter%1000)+1), req)
			case 3:
				s.Delete(fmt.Sprintf("BOOK-%d", (idCounter%1000)+1))
			}
		}
	})
}
