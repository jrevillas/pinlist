package workers

import (
	"time"

	"github.com/go-gorp/gorp"
	"github.com/mvader/pinlist/api/log"
	"github.com/mvader/pinlist/api/models"
)

type Session struct {
	db    *gorp.DbMap
	store models.UserStore
}

func NewSession(db *gorp.DbMap) *Session {
	return &Session{db, models.UserStore{db}}
}

func (s *Session) Run() {
	for _ = range time.Tick(24 * time.Hour) {
		if err := s.store.RemoveExpiredTokens(); err != nil {
			log.Err(err)
		}
	}
}
