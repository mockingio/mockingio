package fixtures

import (
	"github.com/mockingio/engine/mock"
)

func Mock1() *mock.Mock {
	return &mock.Mock{
		ID: "mock1",
		Routes: []*mock.Route{
			{
				ID:     "route1",
				Method: "GET",
				Responses: []mock.Response{
					{
						ID:     "response1",
						Status: 201,
					},
				},
			},
		},
	}
}
