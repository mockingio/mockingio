package engine

import "github.com/mockingio/mockingio/engine/mock"

type Plugin interface {
	Response(response *mock.Response)
}
