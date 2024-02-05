package test

import (
	"snowflake/log"
)

// Start .
func Start() {
	s := NewTester()
	s.Start()
}

// Tester .
type Tester struct {
}

// NewTester .
func NewTester() (s *Tester) {
	s = new(Tester)
	return
}

// Start .
func (s *Tester) Start() {
	log.Infof("?")
}
