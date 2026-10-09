package ports

import (
	"context"
)

// HelloRepo defines the interface (Port) for fetching data from the infrastructure layer.
// Notice it knows nothing about HTTP or the database explicitly.
type HelloRepo interface {
	GetStaticName(ctx context.Context) (string, error)
}
