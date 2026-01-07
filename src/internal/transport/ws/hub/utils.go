package hub

import (
	"encoding/json"

	pkgErrors "github.com/pkg/errors"
)

func unmarshal[T any](data []byte) (*T, error) {
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, pkgErrors.WithStack(err)
	}
	return &v, nil
}
