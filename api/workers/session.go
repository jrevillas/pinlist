package workers

import (
	"time"

	"github.com/mvader/pinlist/api/log"
	"github.com/mvader/pinlist/api/models"
	"gopkg.in/gorp.v1"
)

// Session is a worker in charge of removing expired sessions.
type Session struct {
	db    *gorp.DbMap
	store models.UserStore
}

// NewSession returns a new Session worker.
func NewSession(db *gorp.DbMap) *Session {
	return &Session{db, models.UserStore{DbMap: db}}
}

// Run executes the worker.
func (s *Session) Run() {
	for _ = range time.Tick(24 * time.Hour) {
		if err := s.store.RemoveExpiredTokens(); err != nil {
			log.Err(err)
		}
	}
}
