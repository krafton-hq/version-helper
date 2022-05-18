package build_counter

import (
	"context"
	"fmt"
)

type MemoryCounter struct {
	count uint
}

func NewMemoryCounter(count uint) *MemoryCounter {
	return &MemoryCounter{count: count}
}

func (c *MemoryCounter) String() string {
	return fmt.Sprintf("MemoryCounter: %#v", c)
}

func (c *MemoryCounter) Increase(ctx context.Context) (uint, error) {
	c.count++
	return c.count, nil
}

func (c *MemoryCounter) Get(ctx context.Context) (uint, error) {
	return c.count, nil
}
