package config

import (
	"testing"
)

func TestSetDefaults_ServerPort(t *testing.T) {
	config := &Config{
		Servers: []Server{
			{
				Name:     "server1",
				Host:     "192.168.1.10",
				User:     "admin",
				Password: "secret",
				Port:     0, // Should be set to DefaultSSHPort
			},
			{
				Name:     "server2",
				Host:     "192.168.1.11",
				User:     "admin",
				Password: "secret",
				Port:     2222, // Should remain unchanged
			},
		},
	}

	SetDefaults(config)

	// Check that server1 got the default port
	if config.Servers[0].Port != DefaultSSHPort {
		t.Errorf("expected server1 port to be %d, got %d", DefaultSSHPort, config.Servers[0].Port)
	}

	// Check that server2 port remained unchanged
	if config.Servers[1].Port != 2222 {
		t.Errorf("expected server2 port to remain 2222, got %d", config.Servers[1].Port)
	}
}

func TestSetDefaults_ActionTimeout(t *testing.T) {
	config := &Config{
		Actions: []Action{
			{
				Name:    "action1",
				Command: "echo hello",
				Timeout: 0, // Should be set to DefaultTimeout
			},
			{
				Name:    "action2",
				Command: "echo world",
				Timeout: 60, // Should remain unchanged
			},
		},
	}

	SetDefaults(config)

	// Check that action1 got the default timeout
	if config.Actions[0].Timeout != DefaultTimeout {
		t.Errorf("expected action1 timeout to be %d, got %d", DefaultTimeout, config.Actions[0].Timeout)
	}

	// Check that action2 timeout remained unchanged
	if config.Actions[1].Timeout != 60 {
		t.Errorf("expected action2 timeout to remain 60, got %d", config.Actions[1].Timeout)
	}
}

func TestSetDefaults_MixedConfig(t *testing.T) {
	config := &Config{
		Servers: []Server{
			{
				Name:     "server1",
				Host:     "192.168.1.10",
				User:     "admin",
				Password: "secret",
				Port:     0, // Should be set to DefaultSSHPort
			},
		},
		Actions: []Action{
			{
				Name:    "action1",
				Command: "echo hello",
				Timeout: 0, // Should be set to DefaultTimeout
			},
		},
	}

	SetDefaults(config)

	// Check server defaults
	if config.Servers[0].Port != DefaultSSHPort {
		t.Errorf("expected server port to be %d, got %d", DefaultSSHPort, config.Servers[0].Port)
	}

	// Check action defaults
	if config.Actions[0].Timeout != DefaultTimeout {
		t.Errorf("expected action timeout to be %d, got %d", DefaultTimeout, config.Actions[0].Timeout)
	}
}

func TestSetDefaults_EmptyConfig(t *testing.T) {
	config := &Config{
		Servers: []Server{},
		Actions: []Action{},
	}

	// Should not panic or error
	SetDefaults(config)

	// Config should remain empty
	if len(config.Servers) != 0 {
		t.Error("expected empty servers list to remain empty")
	}
	if len(config.Actions) != 0 {
		t.Error("expected empty actions list to remain empty")
	}
}

func TestDefaultConstants(t *testing.T) {
	// Test that default constants have reasonable values
	if DefaultSSHPort <= 0 {
		t.Errorf("DefaultSSHPort should be positive, got %d", DefaultSSHPort)
	}
	if DefaultSSHPort > 65535 {
		t.Errorf("DefaultSSHPort should be a valid port number, got %d", DefaultSSHPort)
	}

	if DefaultTimeout <= 0 {
		t.Errorf("DefaultTimeout should be positive, got %d", DefaultTimeout)
	}

	if DefaultPasswordLength <= 0 {
		t.Errorf("DefaultPasswordLength should be positive, got %d", DefaultPasswordLength)
	}

	if MaxKeyDirectories <= 0 {
		t.Errorf("MaxKeyDirectories should be positive, got %d", MaxKeyDirectories)
	}
}
