package banek_loader

import "baneks.com/internal/models"

type BaneksLoader interface {
	GetRandomBanek() (models.Banek, error)
}
