package faker

import (
	"github.com/mockingio/mockingio/engine/mock"
)

func New() *Faker {
	return &Faker{}
}

type Faker struct{}

func (f *Faker) Response(response *mock.Response) {

}
