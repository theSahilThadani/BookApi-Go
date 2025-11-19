package store

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"bookapi/internal/model"
)

// ShardedBookStore splits data across N shards to reduce lock contention.
type ShardedBookStore struct {
	shards []*bookShard
	n      uint32
}

type bookShard struct {
	mu    sync.RWMutex
	books map[string]*model.Book
}

func NewShardedBookStore(shardCount uint32) *ShardedBookStore {
	if shardCount == 0 {
		shardCount = 16
	}
	s := &ShardedBookStore{
		n:      shardCount,
		shards: make([]*bookShard, shardCount),
	}
	for i := range s.shards {
		s.shards[i] = &bookShard{books: make(map[string]*model.Book)}
	}
	return s
}

func (s *ShardedBookStore) getShard(key string) *bookShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return s.shards[h.Sum32()%s.n]
}

// Create generates an ID and inserts into a shard
func (s *ShardedBookStore) Create(req model.CreateBookRequest) *model.Book {
	// generate ID based on timestamp+rand pattern for simplicity
	id := fmt.Sprintf("BOOK-%d", time.Now().UnixNano())
	sh := s.getShard(id)
	sh.mu.Lock()
	defer sh.mu.Unlock()

	book := &model.Book{
		ID:          id,
		Title:       req.Title,
		Author:      req.Author,
		ISBN:        req.ISBN,
		PublishYear: req.PublishYear,
		CreatedAt:   time.Now(),
	}
	sh.books[id] = book
	return book
}
func (s *ShardedBookStore) GetAll() []model.Book {
	var result []model.Book
	for _, sh := range s.shards {
		sh.mu.RLock()
		for _, b := range sh.books {
			result = append(result, *b)
		}
		sh.mu.RUnlock()
	}
	return result
}
func (s *ShardedBookStore) GetById(id string) (*model.Book, bool) {
	sh := s.getShard(id)
	sh.mu.RLock()
	defer sh.mu.RUnlock()
	b, ok := sh.books[id]
	return b, ok
}

func (s *ShardedBookStore) Update(id string, req model.UpdateBookRequest) (*model.Book, bool) {
	sh := s.getShard(id)
	sh.mu.Lock()
	defer sh.mu.Unlock()
	b, ok := sh.books[id]
	if !ok {
		return nil, false
	}
	if req.Title != nil {
		b.Title = *req.Title
	}
	if req.Author != nil {
		b.Author = *req.Author
	}
	if req.ISBN != nil {
		b.ISBN = *req.ISBN
	}
	if req.PublishYear != nil {
		b.PublishYear = *req.PublishYear
	}
	return b, true
}

func (s *ShardedBookStore) Delete(id string) bool {
	sh := s.getShard(id)
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if _, ok := sh.books[id]; !ok {
		return false
	}
	delete(sh.books, id)
	return true
}
