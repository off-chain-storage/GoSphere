package runtime

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockService struct {
	status error
}

type secondMockService struct {
	status error
}

func (_ *mockService) Start() {
}

func (_ *mockService) Stop() error {
	return nil
}

func (_ *secondMockService) Start() {
}

func (_ *secondMockService) Stop() error {
	return nil
}

func TestRegisterService_Twice(t *testing.T) {
	registry := &ServiceRegistry{
		services: make(map[reflect.Type]Service),
	}

	m := &mockService{}
	require.NoError(t, registry.RegisterService(m), "Failed to register first mock service")

	require.Equal(t, 1, len(registry.serviceTypes))
	assert.ErrorContains(t, registry.RegisterService(m), "service already exists")
}

func TestRegisterService_Different(t *testing.T) {
	registry := &ServiceRegistry{
		services: make(map[reflect.Type]Service),
	}

	m := &mockService{}
	s := &secondMockService{}
	require.NoError(t, registry.RegisterService(m), "Failed to register first mock service")
	require.NoError(t, registry.RegisterService(s), "Failed to register second mock service")

	require.Equal(t, 2, len(registry.serviceTypes))

	_, exits := registry.services[reflect.TypeOf(m)]
	assert.Equal(t, true, exits, "service of type %v not registered", reflect.TypeOf(m))

	_, exits = registry.services[reflect.TypeOf(m)]
	assert.Equal(t, true, exits, "service of type %v not registered", reflect.TypeOf(s))
}
