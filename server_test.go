package main

import (
	"testing"
	"time"
)

func TestStartServer(t *testing.T) {
	s := &server{}
	timer := time.NewTimer(1 * time.Second)
	go func() {
		<-timer.C
		s.stop()
	}()
	err := s.start("localhost:4208", false, false)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
