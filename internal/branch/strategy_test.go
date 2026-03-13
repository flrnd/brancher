package branch

import "testing"

func TestAvailableStrategies(t *testing.T) {
	list := AvailableStrategies()

	if len(list) == 0 {
		t.Fatalf("expected at least one strategy")
	}
}

func TestResolveStrategy(t *testing.T) {
	s, err := ResolveStrategy(StrategyTitle)
	if err != nil {
		t.Fatalf("ResolveStrategy returned error: %v", err)
	}

	if s == nil {
		t.Fatalf("expected strategy instance, got nil")
	}
}

func TestResolveStrategyUnknown(t *testing.T) {
	_, err := ResolveStrategy("unknown")

	if err == nil {
		t.Fatalf("expected error for unknown strategy")
	}
}
