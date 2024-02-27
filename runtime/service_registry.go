package runtime

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("prefix", "registry")

type Service interface {
	Start()
	Stop() error
}

type ServiceRegistry struct {
	services     map[reflect.Type]Service
	serviceTypes []reflect.Type
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[reflect.Type]Service),
	}
}

func (s *ServiceRegistry) StartAll() {
	log.Debugf("Starting %d services: %v", len(s.serviceTypes), s.serviceTypes)
	for _, kind := range s.serviceTypes {
		log.Debugf("Starting service type %v", kind)
		go s.services[kind].Start()
	}
}

func (s *ServiceRegistry) StopAll() {
	for i := len(s.serviceTypes) - 1; i >= 0; i-- {
		kind := s.serviceTypes[i]
		service := s.services[kind]
		if err := service.Stop(); err != nil {
			log.WithError(err).Errorf("Could not stop the following service: %v", kind)
		}
	}
}

func (s *ServiceRegistry) RegisterService(service Service) error {
	kind := reflect.TypeOf(service)
	if _, exists := s.services[kind]; exists {
		return fmt.Errorf("service already exists: %v", kind)
	}
	s.services[kind] = service
	s.serviceTypes = append(s.serviceTypes, kind)
	return nil
}
