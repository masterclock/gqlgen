package extension

import (
	"context"

	"github.com/99designs/gqlgen/complexity"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/gqlerror"
)

// ComplexityLimit allows you to define a limit on query complexity
//
// If a query is submitted that exceeds the limit, a 422 status code will be returned.
type ComplexityLimit func(ctx context.Context, rc *graphql.OperationContext) int

var _ graphql.OperationContextMutator = ComplexityLimit(func(ctx context.Context, rc *graphql.OperationContext) int { return 0 })

// FixedComplexityLimit sets a complexity limit that does not change
func FixedComplexityLimit(limit int) graphql.HandlerExtension {
	return ComplexityLimit(func(ctx context.Context, rc *graphql.OperationContext) int {
		return limit
	})
}

func (c ComplexityLimit) MutateOperationContext(ctx context.Context, rc *graphql.OperationContext) *gqlerror.Error {
	es := graphql.GetServerContext(ctx)
	op := rc.Doc.Operations.ForName(rc.OperationName)
	complexity := complexity.Calculate(es, op, rc.Variables)

	limit := c(ctx, rc)

	if complexity > limit {
		return gqlerror.Errorf("operation has complexity %d, which exceeds the limit of %d", complexity, limit)
	}

	return nil
}
