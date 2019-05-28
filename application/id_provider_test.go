package application_test

import (
	"testing"

	"github.com/opensaasstudio/meerkat/application"
	"github.com/stretchr/testify/assert"
)

func TestULIDProvider(t *testing.T) {
	p := application.NewULIDProvider()
	assert.NotEqual(t, p.NewString(), p.NewString())
}
