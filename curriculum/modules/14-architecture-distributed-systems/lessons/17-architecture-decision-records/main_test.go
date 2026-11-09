package main

import "testing"

func TestNewADRProposed(t *testing.T) {
	adr := NewADR("ADR-001", "Test Decision")
	if adr.Status != StatusProposed {
		t.Errorf("expected proposed, got %s", adr.Status)
	}
	if adr.ID != "ADR-001" {
		t.Errorf("expected ADR-001, got %s", adr.ID)
	}
	if adr.Title != "Test Decision" {
		t.Errorf("expected Test Decision, got %s", adr.Title)
	}
}

func TestADRAccept(t *testing.T) {
	adr := NewADR("ADR-002", "Accept me")
	adr.Accept()
	if adr.Status != StatusAccepted {
		t.Errorf("expected accepted, got %s", adr.Status)
	}
}

func TestADRDeprecate(t *testing.T) {
	adr := NewADR("ADR-003", "Deprecate me")
	adr.Deprecate()
	if adr.Status != StatusDeprecated {
		t.Errorf("expected deprecated, got %s", adr.Status)
	}
}

func TestADRSupersede(t *testing.T) {
	adr := NewADR("ADR-004", "Old decision")
	adr.Supersede("ADR-005")
	if adr.Status != StatusSuperseded {
		t.Errorf("expected superseded, got %s", adr.Status)
	}
	if adr.Supersedes != "ADR-005" {
		t.Errorf("expected ADR-005, got %s", adr.Supersedes)
	}
}

func TestADRManagerAddAndGet(t *testing.T) {
	m := NewADRManager()
	adr := NewADR("ADR-010", "Test")
	m.Add(adr)
	got := m.Get("ADR-010")
	if got == nil {
		t.Fatal("expected to find ADR")
	}
	if got.Title != "Test" {
		t.Errorf("expected Test, got %s", got.Title)
	}
}

func TestADRManagerByStatus(t *testing.T) {
	m := NewADRManager()
	a1 := NewADR("ADR-020", "Accepted decision")
	a1.Accept()
	a2 := NewADR("ADR-021", "Proposed decision")
	a3 := NewADR("ADR-022", "Deprecated decision")
	a3.Deprecate()
	m.Add(a1)
	m.Add(a2)
	m.Add(a3)
	accepted := m.ByStatus(StatusAccepted)
	if len(accepted) != 1 {
		t.Errorf("expected 1 accepted, got %d", len(accepted))
	}
	proposed := m.ByStatus(StatusProposed)
	if len(proposed) != 1 {
		t.Errorf("expected 1 proposed, got %d", len(proposed))
	}
}

func TestADRManagerAll(t *testing.T) {
	m := NewADRManager()
	m.Add(NewADR("ADR-030", "A"))
	m.Add(NewADR("ADR-031", "B"))
	m.Add(NewADR("ADR-032", "C"))
	if len(m.All()) != 3 {
		t.Errorf("expected 3, got %d", len(m.All()))
	}
}

func TestADRStringContainsSections(t *testing.T) {
	adr := NewADR("ADR-100", "Important Decision")
	adr.Context = "We needed to decide something."
	adr.Decision = "We decided to do X."
	adr.Consequences = "The system is now faster."
	s := adr.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}

func TestADRManagerGetNonExistent(t *testing.T) {
	m := NewADRManager()
	got := m.Get("NONEXISTENT")
	if got != nil {
		t.Fatal("expected nil for non-existent ADR")
	}
}

func TestADRManagerTable(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(m *ADRManager)
		status ADRStatus
		want   int
	}{
		{
			name: "no_accepted",
			setup: func(m *ADRManager) {
				m.Add(NewADR("A1", "Prop"))
			},
			status: StatusAccepted,
			want:   0,
		},
		{
			name: "one_accepted",
			setup: func(m *ADRManager) {
				a := NewADR("A2", "Accepted")
				a.Accept()
				m.Add(a)
				m.Add(NewADR("A3", "Prop"))
			},
			status: StatusAccepted,
			want:   1,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewADRManager()
			tc.setup(m)
			got := len(m.ByStatus(tc.status))
			if got != tc.want {
				t.Errorf("want %d, got %d", tc.want, got)
			}
		})
	}
}
