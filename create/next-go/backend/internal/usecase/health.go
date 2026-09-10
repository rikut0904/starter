package usecase

import "github.com/rikut0904/starter/create/next-go/backend/internal/domain"

type Health struct{}

func NewHealth() Health             { return Health{} }
func (Health) Check() domain.Health { return domain.Health{Status: "ok"} }
