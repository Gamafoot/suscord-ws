package storage

import (
	"suscord_ws/internal/domain/storage/database"
)

type Storage interface {
	Database() database.Storage
}
