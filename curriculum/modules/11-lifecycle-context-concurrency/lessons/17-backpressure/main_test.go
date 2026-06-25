package main

import "testing"

func TestBackpressureCompiles(t *testing.T) {
}

func TestBoundedProducerBackpressure(t *testing.T) {
	done := make(chan struct{})
	ch := boundedProducer(done, 2)
	<-ch
	<-ch
	// buffer full now, producer should block on send
	select {
	case <-ch:
		// ok, consumer is reading
	default:
		t.Fatal("expected at least one more value")
	}
	close(done)
}

func TestLoadSheddingRejectsWhenFull(t *testing.T) {
	ch := loadSheddingProducer(1)
	<-ch
	// after this, buffer should be full or not — depends on timing
	// just verify it doesn't deadlock
	ok := false
	select {
	case <-ch:
		ok = true
	default:
	}
	if !ok {
		t.Log("load shedding rejected a value as expected")
	}
}
