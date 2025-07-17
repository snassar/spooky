package main

import (
	"strings"
	"testing"
)

func TestGenerateID_Unique(t *testing.T) {
	id1 := generateID()
	id2 := generateID()
	if id1 == id2 {
		t.Error("Expected unique IDs, got identical")
	}
	if !strings.Contains(id1, "+") {
		t.Errorf("Expected timestamp+random format, got %q", id1)
	}
}

func TestGenerateGitStyleID_Deterministic(t *testing.T) {
	id1 := generateGitStyleID("foo")
	id2 := generateGitStyleID("foo")
	id3 := generateGitStyleID("bar")
	if id1 != id2 {
		t.Error("Expected deterministic IDs for same input")
	}
	if id1 == id3 {
		t.Error("Expected different IDs for different input")
	}
}

func TestGenerateServers_Count(t *testing.T) {
	scale := ScaleConfig{Name: "test", Hardware: 4, VMs: 8, Description: "test scale"}
	servers := generateServers(scale)
	if len(servers) == 0 {
		t.Errorf("Expected some servers, got %d", len(servers))
	}
}
