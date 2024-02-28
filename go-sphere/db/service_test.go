package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Start_OnlyStartsOnce(t *testing.T) {
	s, _ := setupTest(t)
	require.NotNil(t, s)

	exitRoutine := make(chan bool)
	go func() {
		s.Start()
		<-exitRoutine
	}()
	time.Sleep(time.Second * 2)
	assert.Equal(t, true, s.started, "expected service to be started")
	s.Start()
	require.NoError(t, s.Stop())
	exitRoutine <- true
}

func TestService_Stop_SetsStartedToFalse(t *testing.T) {
	s, _ := setupTest(t)
	require.NotNil(t, s)

	s.SetRedisConn()
	s.started = true
	require.NoError(t, s.Stop())
	assert.Equal(t, false, s.started, "expected service to be stopped")
}
