package main

import (
	"context"
	"fmt"
	"sync"
)

type Tenant string

type contextKey string

const tenantKey contextKey = "tenant"

func WithTenant(ctx context.Context, t Tenant) context.Context {
	return context.WithValue(ctx, tenantKey, t)
}

func TenantFromContext(ctx context.Context) (Tenant, bool) {
	t, ok := ctx.Value(tenantKey).(Tenant)
	return t, ok
}

type TenantStore struct {
	mu   sync.Mutex
	data map[string]map[string]string
}

func NewTenantStore() *TenantStore {
	return &TenantStore{data: make(map[string]map[string]string)}
}

func (s *TenantStore) Set(ctx context.Context, key, value string) error {
	tenant, ok := TenantFromContext(ctx)
	if !ok {
		return fmt.Errorf("no tenant in context")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data[string(tenant)] == nil {
		s.data[string(tenant)] = make(map[string]string)
	}
	s.data[string(tenant)][key] = value
	return nil
}

func (s *TenantStore) Get(ctx context.Context, key string) (string, error) {
	tenant, ok := TenantFromContext(ctx)
	if !ok {
		return "", fmt.Errorf("no tenant in context")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	val, ok := s.data[string(tenant)][key]
	if !ok {
		return "", fmt.Errorf("key not found for tenant %s", tenant)
	}
	return val, nil
}

type TenantConfig struct {
	Tenant          Tenant
	MaxConnections  int
	RateLimitPerSec int
}

type MultiTenantConfig struct {
	mu      sync.RWMutex
	configs map[Tenant]TenantConfig
}

func NewMultiTenantConfig() *MultiTenantConfig {
	return &MultiTenantConfig{configs: make(map[Tenant]TenantConfig)}
}

func (m *MultiTenantConfig) SetConfig(cfg TenantConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.configs[cfg.Tenant] = cfg
}

func (m *MultiTenantConfig) GetConfig(ctx context.Context) (TenantConfig, error) {
	tenant, ok := TenantFromContext(ctx)
	if !ok {
		return TenantConfig{}, fmt.Errorf("no tenant in context")
	}
	m.mu.RLock()
	defer m.mu.RUnlock()
	cfg, ok := m.configs[tenant]
	if !ok {
		return TenantConfig{}, fmt.Errorf("no config for tenant %s", tenant)
	}
	return cfg, nil
}

func (m *MultiTenantConfig) AllConfigs() []TenantConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]TenantConfig, 0, len(m.configs))
	for _, cfg := range m.configs {
		result = append(result, cfg)
	}
	return result
}

func main() {
	ctx1 := WithTenant(context.Background(), Tenant("acme-corp"))
	ctx2 := WithTenant(context.Background(), Tenant("globex-inc"))

	store := NewTenantStore()
	store.Set(ctx1, "api_key", "sk-acme-123")
	store.Set(ctx2, "api_key", "sk-globex-456")

	if v, err := store.Get(ctx1, "api_key"); err == nil {
		fmt.Printf("acme-corp api_key: %s\n", v)
	}
	if v, err := store.Get(ctx2, "api_key"); err == nil {
		fmt.Printf("globex-inc api_key: %s\n", v)
	}

	cfg := NewMultiTenantConfig()
	cfg.SetConfig(TenantConfig{Tenant: "acme-corp", MaxConnections: 10, RateLimitPerSec: 100})
	cfg.SetConfig(TenantConfig{Tenant: "globex-inc", MaxConnections: 5, RateLimitPerSec: 50})

	if c, err := cfg.GetConfig(ctx1); err == nil {
		fmt.Printf("acme-corp config: %+v\n", c)
	}
	if c, err := cfg.GetConfig(ctx2); err == nil {
		fmt.Printf("globex-inc config: %+v\n", c)
	}
}
