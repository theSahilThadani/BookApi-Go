package store

import (
	"bookapi/internal/model"
	"fmt"
	"strconv"
	"sync/atomic"
	"testing"
)

// helper: make a simple create request with i suffix
func mkCreate(i int) model.CreateBookRequest {
	return model.CreateBookRequest{
		Title:       "Title-" + strconv.Itoa(i),
		Author:      "Author",
		ISBN:        fmt.Sprintf("ISBN-%d", i),
		PublishYear: 2024,
	}
}

func BenchmarkStore_CreateParallel(b *testing.B) {
	s := NewBookStore()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			req := mkCreate(i)
			s.Create(req)
			i++
		}
	})
}

func BenchmarkStore_MixedParallel(b *testing.B) {
	s := NewBookStore()
	var idCounter uint64

	// pre-create some entries
	for i := 0; i < 1000; i++ {
		s.Create(mkCreate(i))
		atomic.AddUint64(&idCounter, 1)
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {

			n := int(atomic.AddUint64(&idCounter, 1) % 4)
			switch n {
			case 0:
				s.Create(mkCreate(int(idCounter)))
			case 1:

				s.GetById(fmt.Sprintf("BOOK-%d", (idCounter%1000)+1))
			case 2:

				req := model.UpdateBookRequest{Title: strPtr("Updated")}
				s.Update(fmt.Sprintf("BOOK-%d", (idCounter%1000)+1), req)
			case 3:
				s.Delete(fmt.Sprintf("BOOK-%d", (idCounter%1000)+1))
			}
		}
	})
}

func strPtr(s string) *string { return &s }
