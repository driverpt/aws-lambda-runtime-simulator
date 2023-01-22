package event

import (
	"errors"
	"lambda-runtime-simulator/pkg/config"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Invocation struct {
	Id        string        `json:"id"`
	Body      string        `json:"body"`
	Timeout   time.Time     `json:"timeout"`
	Response  *string       `json:"response"`
	Error     *RuntimeError `json:"error"`
	ErrorType *string       `json:"errorType"`
}

type Service struct {
	config  *config.Runtime
	holder  map[string]*Invocation
	channel chan *Invocation
}

func NewService(cfg *config.Runtime) *Service {
	result := &Service{
		config:  cfg,
		holder:  map[string]*Invocation{},
		channel: make(chan *Invocation, 100),
	}

	return result
}

func (s Service) ResetAll() error {
	log.Warn("Resetting internal cache")
	prevChan := s.channel
	s.channel = make(chan *Invocation, 100)
	if prevChan != nil {
		close(prevChan)
	}

	s.holder = map[string]*Invocation{}

	return nil
}

func (s Service) GetNextInvocation() (*Invocation, error) {
	log.Debug("Awaiting next invocation")
	next := <-s.channel
	// TODO: Some kind of error handling here
	log.Debugf("Next invocation received %v", next.Id)
	return next, nil
}

func (s Service) PushInvocation(body string) (string, error) {
	invocation := Invocation{
		Id:      uuid.NewString(),
		Body:    body,
		Timeout: time.Now().UTC().Add(time.Duration(s.config.TimeoutInSeconds) * time.Second),
	}

	s.holder[invocation.Id] = &invocation
	s.channel <- &invocation

	return invocation.Id, nil
}

func (s Service) SendResponse(id string, body []byte) error {
	inv := s.holder[id]
	if inv == nil {
		return errors.New("request does not exist")
	}

	now := time.Now().UTC()
	// Check for Timeout
	if now.After(inv.Timeout) || now.Equal(inv.Timeout) {
		return errors.New("invocation timeout")
	}

	b := string(body)
	inv.Body = b
	return nil
}

func (s Service) SendError(id string, error *RuntimeError, errorType string) error {
	inv := s.holder[id]
	if inv == nil {
		return errors.New("request does not exist")
	}

	now := time.Now().UTC()
	// Check for Timeout
	if now.After(inv.Timeout) || now.Equal(inv.Timeout) {
		return errors.New("invocation timeout")
	}

	inv.Error = error
	if errorType != "" {
		inv.ErrorType = &errorType
	}

	return nil
}

func (s Service) GetCachedInvocations() []*Invocation {
	var result []*Invocation

	for _, v := range s.holder {
		result = append(result, v)
	}

	return result
}

func (s Service) GetById(id string) *Invocation {
	return s.holder[id]
}
