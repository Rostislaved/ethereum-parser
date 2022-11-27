package inMemoryStorage

import (
	"github.com/Rostislaved/ethereum-parser/internal/app/entity"
	"sync"
)

type Storage struct {
	m    map[string][]entity.Transaction
	rwmu *sync.RWMutex
}

func New() *Storage {
	m := make(map[string][]entity.Transaction)

	rwmu := sync.RWMutex{}

	return &Storage{
		m:    m,
		rwmu: &rwmu,
	}
}
