package cutofftimes

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func makeCutoffTimesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.parseInwardCutttoffTimes(req)
		return nil, nil
	}
}
