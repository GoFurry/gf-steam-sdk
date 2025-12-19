package store

import (
	"fmt"

	"github.com/GoFurry/gf-steam-sdk/internal/client"
)

type StoreService struct {
	client *client.Client
}

func NewStoreService(c *client.Client) *StoreService {
	return &StoreService{client: c}
}

// Close 释放StoreService资源
func (s *StoreService) Close() error {
	if s.client == nil {
		return nil
	}
	if err := s.client.Close(); err != nil {
		return fmt.Errorf("StoreService Close failed: %w", err)
	}
	return nil
}
