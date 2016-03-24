package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ipEntry struct {
	ip   string
	hits int
	last time.Time
}

type ipRegister struct {
	mut     sync.RWMutex
	entries map[string]ipEntry
}

var (
	register         = &ipRegister{entries: make(map[string]ipEntry)}
	cleanupFrequency = 15 * time.Minute
)

const maxHits = 20

func GuestHourLimit(c *gin.Context) {
	addr := getIP(c.Request)

	register.mut.Lock()
	if v, ok := register.entries[addr]; ok && v.hits > maxHits {
		c.String(http.StatusForbidden, "this ip address is blocked for an hour")
		register.mut.Unlock()
		return
	} else if ok {
		v.hits++
		v.last = time.Now()
		register.entries[addr] = v
	} else {
		register.entries[addr] = ipEntry{
			ip:   addr,
			hits: 1,
			last: time.Now(),
		}
	}

	register.mut.Unlock()
	c.Next()
}

func getIP(r *http.Request) string {
	addr := r.Header.Get("X-Real-IP")
	if addr == "" {
		addr = r.Header.Get("X-Forwarded-For")
		if addr == "" {
			addr = r.RemoteAddr
		}
	}
	return addr
}

func cleanRegister() {
	for _ = range time.Tick(cleanupFrequency) {
		register.mut.Lock()
		for k, v := range register.entries {
			if v.last.Add(1 * time.Hour).Before(time.Now()) {
				delete(register.entries, k)
			}
		}
		register.mut.Unlock()
	}
}

func main() {
	go cleanRegister()
}
