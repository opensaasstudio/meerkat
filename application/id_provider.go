package application

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid"
	"github.com/opensaasstudio/meerkat/domain"
)

type ULIDProvider struct {
	entropy io.Reader
}

func NewULIDProvider() ULIDProvider {
	return ULIDProvider{
		entropy: ulid.Monotonic(rand.New(newPooledRandSource()), 0),
	}
}

func (p ULIDProvider) NewString() string {
	return ulid.MustNew(ulid.Now(), p.entropy).String()
}

func (p ULIDProvider) NewAnswererID() domain.AnswererID {
	return domain.AnswererID(p.NewString())
}

func (p ULIDProvider) NewQuestionnaireID() domain.QuestionnaireID {
	return domain.QuestionnaireID(p.NewString())
}

func (p ULIDProvider) NewQuestionID() domain.QuestionID {
	return domain.QuestionID(p.NewString())
}

func (p ULIDProvider) NewAnswerID() domain.AnswerID {
	return domain.AnswerID(p.NewString())
}

func (p ULIDProvider) NewNotificationTargetID() domain.NotificationTargetID {
	return domain.NotificationTargetID(p.NewString())
}

type pooledRandSource struct {
	pool *sync.Pool
}

func newPooledRandSource() rand.Source {
	return pooledRandSource{
		pool: &sync.Pool{
			New: func() interface{} {
				return rand.NewSource(time.Now().UnixNano())
			},
		},
	}
}

func (s pooledRandSource) Int63() int64 {
	rs := s.pool.Get().(rand.Source)
	defer s.pool.Put(rs)
	return rs.Int63()
}

func (s pooledRandSource) Seed(seed int64) {
	// noop
}
