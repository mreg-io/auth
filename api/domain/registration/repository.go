package registration

import "context"

type Repository interface {
	CreateFlow(ctx context.Context, flow *Flow) error
	QueryFlowByFlowID(ctx context.Context, flow *Flow) error
}
