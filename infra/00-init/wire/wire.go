//go:build wireinjection
// +build wireinjection

package wire

import (
	"PetTrack/infra/00-core/model"

	"github.com/google/wire"
)

func InitService(cfg model.Config) (*Services, error) {
	wire.Build(
		machineSet,
		repoSet,
		serviceSet,
	)
	return &Services{}, nil
}
