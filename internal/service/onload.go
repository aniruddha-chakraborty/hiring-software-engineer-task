package service

import (
	"go.uber.org/zap"
)

type Onload struct {
	log *zap.SugaredLogger
	c   *Cache
}

func NewOnloadService(log *zap.SugaredLogger, c *Cache) *Onload {
	return &Onload{
		log: log,
		c:   c,
	}
}

func (o *Onload) Start() {
	o.c.PopulateCache()
}
