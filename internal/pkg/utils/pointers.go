package pkg

import "github.com/google/uuid"

func PointerTo[T any](v T) *T {
	return &v
}

func UUID(id string) *uuid.UUID {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil
	}
	return &uid
}
