package registration

import "context"

type Repository interface {
	CreateFlow(ctx context.Context, flow *Flow) error
	QueryFlow(ctx context.Context, flow *Flow) error
}
