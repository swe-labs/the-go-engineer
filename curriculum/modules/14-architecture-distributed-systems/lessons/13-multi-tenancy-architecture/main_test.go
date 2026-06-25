package main

import (
	"context"
	"testing"
)

func TestWithTenantAndExtract(t *testing.T) {
	ctx := WithTenant(context.Background(), Tenant("my-tenant"))
	got, ok := TenantFromContext(ctx)
	if !ok {
		t.Fatal("expected tenant in context")
	}
	if got != "my-tenant" {
		t.Errorf("expected my-tenant, got %s", got)
	}
}

func TestTenantFromContextMissing(t *testing.T) {
	_, ok := TenantFromContext(context.Background())
	if ok {
		t.Fatal("expected no tenant in bare context")
	}
}

func TestTenantStoreIsolation(t *testing.T) {
	store := NewTenantStore()
	ctxA := WithTenant(context.Background(), Tenant("tenant-a"))
	ctxB := WithTenant(context.Background(), Tenant("tenant-b"))

	store.Set(ctxA, "key1", "value-a")
	store.Set(ctxB, "key1", "value-b")

	gotA, err := store.Get(ctxA, "key1")
	if err != nil {
		t.Fatal(err)
	}
	if gotA != "value-a" {
		t.Errorf("tenant-a expected value-a, got %s", gotA)
	}

	gotB, err := store.Get(ctxB, "key1")
	if err != nil {
		t.Fatal(err)
	}
	if gotB != "value-b" {
		t.Errorf("tenant-b expected value-b, got %s", gotB)
	}
}

func TestTenantStoreMissingKey(t *testing.T) {
	store := NewTenantStore()
	ctx := WithTenant(context.Background(), Tenant("t"))
	_, err := store.Get(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestTenantStoreMissingTenant(t *testing.T) {
	store := NewTenantStore()
	err := store.Set(context.Background(), "k", "v")
	if err == nil {
		t.Fatal("expected error on set without tenant")
	}
}

func TestMultiTenantConfigIsolation(t *testing.T) {
	m := NewMultiTenantConfig()
	m.SetConfig(TenantConfig{Tenant: "tenant-a", MaxConnections: 10})
	m.SetConfig(TenantConfig{Tenant: "tenant-b", MaxConnections: 20})

	ctxA := WithTenant(context.Background(), Tenant("tenant-a"))
	cfg, err := m.GetConfig(ctxA)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.MaxConnections != 10 {
		t.Errorf("expected 10, got %d", cfg.MaxConnections)
	}
}

func TestMultiTenantConfigAll(t *testing.T) {
	m := NewMultiTenantConfig()
	m.SetConfig(TenantConfig{Tenant: "a"})
	m.SetConfig(TenantConfig{Tenant: "b"})
	m.SetConfig(TenantConfig{Tenant: "c"})
	all := m.AllConfigs()
	if len(all) != 3 {
		t.Errorf("expected 3 configs, got %d", len(all))
	}
}

func TestMultiTenantConfigTable(t *testing.T) {
	tests := []struct {
		name   string
		tenant Tenant
		maxCon int
		wantOK bool
	}{
		{name: "configured_tenant", tenant: "t1", maxCon: 10, wantOK: true},
		{name: "another_tenant", tenant: "t2", maxCon: 5, wantOK: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMultiTenantConfig()
			m.SetConfig(TenantConfig{Tenant: tc.tenant, MaxConnections: tc.maxCon})
			ctx := WithTenant(context.Background(), tc.tenant)
			cfg, err := m.GetConfig(ctx)
			if tc.wantOK && err != nil {
				t.Fatal(err)
			}
			if tc.wantOK && cfg.MaxConnections != tc.maxCon {
				t.Errorf("expected %d, got %d", tc.maxCon, cfg.MaxConnections)
			}
		})
	}
}
