package application

//go:generate mockgen -source $GOFILE -destination mock_$GOPACKAGE/${GOFILE}_mock.go

import (
	"context"

	"github.com/opensaasstudio/meerkat/domain"
)

//genconstructor
type CreatingAnswererUsecase struct {
	repository         AnswererRepository `required:""`
	answererIDProvider AnswererIDProvider `required:""`
}

//genconstructor
type CreatingAnswererUsecaseInput struct {
	name string `required:"" getter:""`
}

func (u CreatingAnswererUsecase) CreateAnswerer(ctx context.Context, input CreatingAnswererUsecaseInput) (domain.Answerer, domain.Error) {
	answerer := domain.NewAnswerer(
		u.answererIDProvider.NewAnswererID(),
		input.Name(),
	)
	err := u.repository.Create(ctx, answerer)
	return answerer, err
}

type AnswererIDProvider interface {
	NewAnswererID() domain.AnswererID
}
