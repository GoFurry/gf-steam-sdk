package dev

import (
	"fmt"

	"github.com/GoFurry/gf-steam-sdk/internal/client"
)

type DevService struct {
	client *client.Client
}

func NewDevService(c *client.Client) *DevService {
	return &DevService{client: c}
}

func (g *DevService) Close() error {
	if g.client == nil {
		return nil
	}
	if err := g.client.Close(); err != nil {
		return fmt.Errorf("DevService Close failed: %w", err)
	}
	return nil
}
