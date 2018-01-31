package cutofftimes

import (
	"github.com/go-kit/kit/endpoint"
	"context"
)

func makeCutoffTimesEndpoint(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.parseInwardCutttoffTimes(req)
		return nil, nil
	}
}
