package memory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/mockingio/mockingio/engine/mock"
	"github.com/mockingio/mockingio/engine/persistent"
)

var _ persistent.Persistent = &Memory{}

type MemoryRoute struct {
	mockID string
	route  *mock.Route
}

type MemoryResponse struct {
	mockID   string
	routeID  string
	response *mock.Response
}

type MemoryRule struct {
	mockID     string
	routeID    string
	responseID string
	rule       *mock.Rule
}

type Memory struct {
	mu          sync.Mutex
	configs     map[string]*mock.Mock
	routes      map[string]MemoryRoute
	responses   map[string]MemoryResponse
	rules       map[string]MemoryRule
	kv          map[string]any
	subscribers []func(mock mock.Mock)
}

func New() *Memory {
	return &Memory{
		configs: map[string]*mock.Mock{},
		kv:      map[string]any{},
	}
}

func (m *Memory) SubscribeMockChanges(subscriber func(mock mock.Mock)) {
	m.subscribers = append(m.subscribers, subscriber)
}

func (m *Memory) Get(_ context.Context, key string) (any, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.kv[key], nil
}

func (m *Memory) Set(_ context.Context, key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.kv[key] = value
	return nil
}

func (m *Memory) SetMock(_ context.Context, cfg *mock.Mock) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.configs[cfg.ID] = cfg
	m.routes = map[string]MemoryRoute{}
	m.responses = map[string]MemoryResponse{}
	m.rules = map[string]MemoryRule{}
	for _, route := range cfg.Routes {
		m.routes[route.ID] = MemoryRoute{
			mockID: cfg.ID,
			route:  route,
		}
		for _, response := range route.Responses {
			m.responses[response.ID] = MemoryResponse{
				mockID:   cfg.ID,
				routeID:  route.ID,
				response: &response,
			}

			for _, rule := range response.Rules {
				r := rule
				m.rules[rule.ID] = MemoryRule{
					mockID:     cfg.ID,
					routeID:    route.ID,
					responseID: response.ID,
					rule:       &r,
				}
			}
		}
	}

	for _, subscriber := range m.subscribers {
		subscriber(*cfg)
	}

	return nil
}

func (m *Memory) SaveMock(ctx context.Context, id string) error {
	cfg, error := m.GetMock(ctx, id)

	if error != nil {
		return nil
	}

	for _, subscriber := range m.subscribers {
		subscriber(*cfg)
	}

	return nil
}

func (m *Memory) GetMock(_ context.Context, id string) (*mock.Mock, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cfg, ok := m.configs[id]

	if !ok {
		return nil, nil
	}

	cfg.Routes = []*mock.Route{}

	// FIXME: Avoid to using multiple loop
	for _, route := range m.routes {
		if route.mockID == cfg.ID {
			route.route.Responses = []mock.Response{}
			for _, response := range m.responses {
				if response.routeID == route.route.ID {
					response.response.Rules = []mock.Rule{}
					for _, rule := range m.rules {
						if rule.responseID == response.response.ID {
							response.response.Rules = append(response.response.Rules, *rule.rule)
						}
					}
					route.route.Responses = append(route.route.Responses, *response.response)
				}
			}
			cfg.Routes = append(cfg.Routes, route.route)
		}
	}

	return cfg, nil
}

func (m *Memory) GetMocks(_ context.Context) ([]*mock.Mock, error) {
	var configs []*mock.Mock
	for _, cfg := range m.configs {
		configs = append(configs, cfg)
	}

	return configs, nil
}

func (m *Memory) GetInt(ctx context.Context, key string) (int, error) {
	v, err := m.Get(ctx, key)
	if err != nil {
		return 0, err
	}

	if v == nil {
		return 0, nil
	}

	value, ok := v.(int)
	if !ok {
		return 0, nil
	}

	return value, nil
}

func (m *Memory) Increment(_ context.Context, key string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, ok := m.kv[key]
	if !ok {
		m.kv[key] = 1
		return 1, nil
	}

	val, ok := value.(int)
	if !ok {
		return 0, errors.New(fmt.Sprintf("unable to increase non-int key (%s)", key))
	}

	val++
	m.kv[key] = val

	return val, nil
}

func (m *Memory) SetActiveSession(ctx context.Context, mockID string, sessionID string) error {
	return m.Set(ctx, toActiveSessionKey(mockID), sessionID)
}

func (m *Memory) GetActiveSession(ctx context.Context, mockID string) (string, error) {
	value, err := m.Get(ctx, toActiveSessionKey(mockID))
	if err != nil {
		return "", err
	}

	if v, ok := value.(string); ok {
		return v, nil
	}

	return "", errors.New("unable to convert to string value")
}

func (m *Memory) GetRoute(ctx context.Context, mockID string, routeID string) (*mock.Route, error) {
	_, exists := m.configs[mockID]

	if !exists {
		return nil, errors.New("Mock not found")
	}

	route, exists := m.routes[routeID]

	if !exists {
		return nil, errors.New("route not found")
	}

	return route.route, nil
}

func (m *Memory) PatchRoute(ctx context.Context, mockID string, routeID string, data string) error {
	route, err := m.GetRoute(ctx, mockID, routeID)

	if err != nil {
		return err
	}

	var values map[string]*json.RawMessage
	if err := json.Unmarshal([]byte(data), &values); err != nil {
		return err
	}

	if err := patchStruct(route, values); err != nil {
		return err
	}

	if err := m.SaveMock(ctx, mockID); err != nil {
		return err
	}

	return nil
}

func (m *Memory) DeleteRoute(ctx context.Context, mockID string, routeID string) error {
	_, err := m.GetRoute(ctx, mockID, routeID)

	if err != nil {
		return err
	}

	delete(m.routes, routeID)

	if err := m.SaveMock(ctx, mockID); err != nil {
		return err
	}

	return nil
}

func (m *Memory) CreateRoute(ctx context.Context, mockID string, newRoute mock.Route) error {
	mok, err := m.GetMock(ctx, mockID)
	if err != nil {
		return err
	}

	if mok == nil {
		return errors.New("mock not found")
	}

	_, err = m.GetRoute(ctx, mockID, newRoute.ID)

	if err == nil {
		return errors.New("route already created")
	}

	mok.Routes = append(mok.Routes, &newRoute)

	if err := m.SetMock(ctx, mok); err != nil {
		return err
	}

	return nil
}

func (m *Memory) GetResponse(ctx context.Context, mockID string, responseID string) (*mock.Response, error) {
	mock, exists := m.configs[mockID]

	if !exists || mock == nil {
		return nil, errors.New("Mock not found")
	}

	response, exists := m.responses[responseID]

	if !exists {
		return nil, errors.New("route not found")
	}

	return response.response, nil
}

func (m *Memory) PatchResponse(ctx context.Context, mockID, responseID, data string) error {
	response, err := m.GetResponse(ctx, mockID, responseID)
	if err != nil {
		return err
	}

	var values map[string]*json.RawMessage
	if err := json.Unmarshal([]byte(data), &values); err != nil {
		return err
	}

	if err := patchStruct(response, values); err != nil {
		return err
	}

	if err := m.SaveMock(ctx, mockID); err != nil {
		return err
	}

	return nil
}

func (m *Memory) GetRule(ctx context.Context, mockID string, ruleID string) (*mock.Rule, error) {
	_, exists := m.configs[mockID]

	if !exists {
		return nil, errors.New("Mock not found")
	}

	rule, exists := m.rules[ruleID]

	if !exists {
		return nil, errors.New("route not found")
	}

	return rule.rule, nil
}

func (m *Memory) CreateRule(ctx context.Context, mockID string, responseID string, newRule mock.Rule) error {
	_, err := m.GetResponse(ctx, mockID, responseID)
	if err != nil {
		return err
	}

	_, err = m.GetRule(ctx, mockID, newRule.ID)

	if err == nil {
		return errors.New("Rule already created")
	}

	m.rules[newRule.ID] = MemoryRule{
		mockID:     mockID,
		responseID: responseID,
		rule:       &newRule,
	}

	if err := m.SaveMock(ctx, mockID); err != nil {
		return err
	}

	return nil
}

func (m *Memory) DeleteRule(ctx context.Context, mockID string, ruleID string) error {
	_, err := m.GetRule(ctx, mockID, ruleID)

	if err != nil {
		return errors.New("Rule not found")
	}

	delete(m.rules, ruleID)

	if err := m.SaveMock(ctx, mockID); err != nil {
		return err
	}

	return nil
}

func toActiveSessionKey(mockID string) string {
	return fmt.Sprintf("%s-active-session", mockID)
}

func patchStruct(resource interface{}, patches map[string]*json.RawMessage) error {
	value := reflect.ValueOf(resource)
	for value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("can't operate on non-struct: %s", value.Kind().String())
	}
	if !value.CanAddr() {
		return errors.New("unaddressable struct value")
	}
	valueT := value.Type()
	for i := 0; i < valueT.NumField(); i++ {
		field := value.Field(i)
		if !field.CanAddr() || !field.CanInterface() {
			continue
		}
		if patch, ok := patches[jsonFieldName(valueT.Field(i))]; ok {
			field.Set(reflect.Zero(field.Type()))
			if err := json.Unmarshal(*patch, field.Addr().Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}

func jsonFieldName(field reflect.StructField) string {
	name := strings.Split(field.Tag.Get("json"), ",")[0]
	if name == "" {
		name = field.Name
	}
	return name
}
