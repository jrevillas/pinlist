package workers

import (
	"time"

	"github.com/go-gorp/gorp"
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
	tick := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-tick.C:
			if err := s.store.RemoveExpiredTokens(); err != nil {
				// TODO: Log error
			}
		case <-time.After(5 * time.Second):
		}
	}
}
