package store

import (
	"github.com/GoFurry/gf-steam-sdk/internal/client"
)

type StoreService struct {
	client *client.Client
}

func NewStoreService(c *client.Client) *StoreService {
	return &StoreService{client: c}
}
