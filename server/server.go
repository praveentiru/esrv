package server

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/praveentiru/efp"
)

const ttl = 30

type cItem struct {
	eval       gval.Evaluable
	ttl        time.Duration
	lastAccess time.Time
}

// Server struct holds state of server
type Server struct {
	cache map[string]cItem
	cStop chan int
	mux   sync.Mutex
}

// New allows creation of a new server instance
func New() *Server {
	// Init cleanup routine here
	cMap := make(map[string]cItem)
	cs := make(chan int)
	srv := Server{
		cache: cMap,
		cStop: cs,
	}
	go srv.cClean()
	return &srv
}

// EvalExpression evaluates EXCEL equation and returns value
func (s *Server) EvalExpression(exp, outType string, params map[string]interface{}) (interface{}, error) {
	eval, err := s.getEvaluable(exp)
	if err != nil {
		return nil, err
	}
	switch outType {
	case "string":
		return eval.EvalString(context.Background(), params)
	case "int":
		return eval.EvalInt(context.Background(), params)
	case "boolean":
		return eval.EvalBool(context.Background(), params)
	}
	return nil, errors.New("Request for unsupported return value type")
}

// Stop triggers stop signal for cache cleanup go routine
func (s *Server) Stop() {
	s.cStop <- 1
	return
}

func (s *Server) getEvaluable(exp string) (gval.Evaluable, error) {
	r := s.lookupCache(exp)
	if r != nil {
		return r, nil
	}
	r, err := parse(exp)
	if err != nil {
		return nil, err
	}
	ci := buildCacheItem(r)
	s.addKey(exp, ci)
	return r, nil
}

func (s *Server) lookupCache(exp string) gval.Evaluable {
	c, ok := s.cache[exp]
	if !ok {
		return nil
	}
	c.lastAccess = time.Now()
	return c.eval
}

// cClean is cache manager where it checks and removes items that have become stale.
// Items are marked stale if they have not be called on for (ttl) seconds
func (s *Server) cClean() {
	for {
		select {
		case <-s.cStop:
			break
		case t := <-time.After(time.Duration(ttl) * time.Second):
			for k, v := range s.cache {
				tl := t.Sub(v.lastAccess)
				if tl > time.Duration(ttl)*time.Second {
					s.removeKey(k)
				}
			}
		}
	}
}

func (s *Server) addKey(k string, v cItem) {
	s.mux.Lock()
	s.cache[k] = v
	s.mux.Unlock()
}

func (s *Server) removeKey(k string) {
	s.mux.Lock()
	delete(s.cache, k)
	s.mux.Unlock()
}

func parse(exp string) (gval.Evaluable, error) {
	rdr := strings.NewReader(exp)
	eval, err := efp.Parse(rdr)
	if err != nil {
		return nil, err
	}
	return eval, nil
}

func buildCacheItem(eval gval.Evaluable) cItem {
	tl := time.Duration(ttl) * time.Second
	r := cItem{
		eval:       eval,
		ttl:        tl,
		lastAccess: time.Now()}
	return r
}
