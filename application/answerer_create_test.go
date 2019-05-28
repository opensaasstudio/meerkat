package application_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/opensaasstudio/meerkat/application"
	"github.com/opensaasstudio/meerkat/application/mock_application"
	"github.com/opensaasstudio/meerkat/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreatingAnswererUsecase(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repository := mock_application.NewMockAnswererRepository(ctrl)
	answererIDProvider := mock_application.NewMockAnswererIDProvider(ctrl)
	u := application.NewCreatingAnswererUsecase(repository, answererIDProvider)

	input := application.NewCreatingAnswererUsecaseInput("name")

	answererIDProvider.EXPECT().NewAnswererID().Return(domain.AnswererID("id"))

	repository.EXPECT().Create(ctx, gomock.Any()).Return(nil)

	answerer, err := u.CreateAnswerer(ctx, input)
	assert.NoError(t, err)
	assert.Equal(t, domain.AnswererID("id"), answerer.ID())
}
