// Copyright 2019 The meerkat Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
