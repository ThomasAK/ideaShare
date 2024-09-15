package testutil

import (
	"github.com/go-faker/faker/v4"
	"ideashare/models"
)

func MakeFake[T models.BaseModel](obj T) T {
	err := faker.FakeData(obj)
	if err != nil {
		panic(err)
	}
	return obj
}
