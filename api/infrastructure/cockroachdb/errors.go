package cockroachdb

import "errors"

var ErrConstraint = errors.New("input values failed to meet schema constraint")
