package app

import (
	"suscord_ws/internal/config"
	"suscord_ws/internal/domain/storage"
	"suscord_ws/internal/domain/storage/database"
	"suscord_ws/internal/domain/storage/file"
	"suscord_ws/internal/infrastructure/database/relational"

	pkgErrors "github.com/pkg/errors"
)

type storageImpl struct {
	dbStorage database.Storage
	file      file.FileStorage
}

func NewStorage(cfg *config.Config) (storage.Storage, error) {
	dbStorage, err := relational.NewStorage(cfg.Database.Addr, cfg.Database.LogLevel)
	if err != nil {
		return nil, pkgErrors.WithStack(err)
	}

	return &storageImpl{
		dbStorage: dbStorage,
	}, nil
}

func (s *storageImpl) Database() database.Storage {
	return s.dbStorage
}
