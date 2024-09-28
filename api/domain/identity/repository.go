package identity

import "context"

type Repository interface {
	CreateIdentity(ctx context.Context, identity *Identity) error
	QueryEmail(ctx context.Context, email *Email) error
	EmailExists(ctx context.Context, email string) (bool, error)
}
