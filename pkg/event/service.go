package event

import (
	"encoding/json"
	"errors"
	"lambda-runtime-simulator/pkg/config"
	"lambda-runtime-simulator/pkg/utils"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Invocation struct {
	Id        string        `json:"id"`
	Body      interface{}   `json:"body,omitempty"`
	Timeout   time.Time     `json:"timeout"`
	Response  interface{}   `json:"response,omitempty"`
	StartedAt time.Time     `json:"startedAt"`
	Duration  *int          `json:"duration,omitempty"`
	Error     *RuntimeError `json:"error,omitempty"`
	ErrorType *string       `json:"errorType,omitempty"`
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

func (s *Service) ResetAll() error {
	log.Warn("Resetting internal cache. All pending invocations will return nil")
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
	if next != nil {
		log.Debugf("Next invocation received %v", next.Id)
	}
	return next, nil
}

func (s Service) PushInvocation(body string) (string, error) {
	log.Debugf("Pushing new invocation: %v", body)
	var b interface{}
	err := json.Unmarshal([]byte(body), &b)
	if err != nil {
		log.Errorf("Invalid JSON: %v", err)
		return "", err
	}

	now := time.Now().UTC()

	invocation := Invocation{
		Id:        uuid.NewString(),
		Body:      b,
		StartedAt: now,
	}

	invocation.Timeout = now.Add(time.Duration(s.config.TimeoutInSeconds) * time.Second)
	log.Debugf("Invocation ID: %v", invocation.Id)

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

	var b interface{}
	err := json.Unmarshal(body, &b)
	if err != nil {
		return err
	}

	log.Debugf("Response received %v: %v", id, b)
	diff := now.Sub(inv.StartedAt)

	inv.Response = b
	inv.Duration = utils.ToPointer(int(diff.Seconds()))
	return nil
}

func (s Service) SendError(id string, error *RuntimeError, errorType string) error {
	inv := s.holder[id]
	if inv == nil {
		return errors.New("request does not exist")
	}

	if inv.Response != nil || inv.Error != nil {
		return errors.New("response already set")
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

	log.Debugf("Error received %v: %v,%v", id, error, errorType)

	diff := now.Sub(inv.StartedAt)
	inv.Duration = utils.ToPointer(int(diff.Seconds()))
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

func (s Service) SendRuntimeInitError(message string, errorType string, stackTrace ...string) error {
	// For now just Log the request
	log.Infof("Runtime Error Received [Message:%v, Type:%v]", message, errorType)
	log.Infof("StackTrace: %v", stackTrace)
	return nil
}
