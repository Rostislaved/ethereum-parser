package inMemoryStorage

import (
	"github.com/Rostislaved/ethereum-parser/internal/app/entity"
)

func (s *Storage) SaveTransaction(address string, transaction entity.Transaction) error {
	s.rwmu.Lock()
	defer s.rwmu.Unlock()

	s.m[address] = append(s.m[address], transaction)

	return nil
}

func (s *Storage) GetTransactions(address string) (transactions []entity.Transaction, err error) {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	transactions, ok := s.m[address]
	if !ok {
		return nil, nil
	}

	return transactions, nil
}

func (s *Storage) GetStorageInfo() (entity.StorageInfo, error) {
	s.rwmu.RLock()
	defer s.rwmu.RUnlock()

	numberOfAddressesTransactions := make(map[string]int)

	for key, values := range s.m {
		numberOfAddressesTransactions[key] = len(values)
	}

	storageInfo := entity.StorageInfo{
		NumberOfAddresses:             len(s.m),
		NumberOfAddressesTransactions: numberOfAddressesTransactions,
	}

	return storageInfo, nil
}
