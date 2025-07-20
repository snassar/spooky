package facts

import (
	"context"
	"fmt"
	"net"
	"os"
	"spooky/internal/config"
	"spooky/internal/ssh"
	"testing"
	"time"

	gliderssh "github.com/gliderlabs/ssh"
)

// MockSSHServer represents a mock SSH server for testing
type MockSSHServer struct {
	server *gliderssh.Server
	port   int
	addr   string
}

// NewMockSSHServer creates a new mock SSH server
func NewMockSSHServer() (*MockSSHServer, error) {
	// Find an available port on localhost only
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("failed to find available port: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// Create SSH server
	s := &gliderssh.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", port),
		Handler: func(s gliderssh.Session) {
			// Handle different commands
			cmd := s.Command()
			if len(cmd) == 0 {
				if _, err := s.Write([]byte("mock-ssh-server")); err != nil {
					return
				}
				return
			}

			// Handle commands with two arguments
			if len(cmd) >= 2 {
				switch cmd[0] + " " + cmd[1] {
				case "hostname -f":
					if _, err := s.Write([]byte("test-host.example.com\n")); err != nil {
						return
					}
					return
				case "cat /etc/machine-id":
					if _, err := s.Write([]byte("test-machine-id-12345\n")); err != nil {
						return
					}
					return
				case "cat /etc/os-release":
					if _, err := s.Write([]byte(`NAME="Ubuntu"
VERSION="22.04.3 LTS (Jammy Jellyfish)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 22.04.3 LTS"
VERSION_ID="22.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
UBUNTU_CODENAME=jammy
`)); err != nil {
						return
					}
					return
				case "uname -m":
					if _, err := s.Write([]byte("x86_64\n")); err != nil {
						return
					}
					return
				case "uname -r":
					if _, err := s.Write([]byte("5.15.0-88-generic\n")); err != nil {
						return
					}
					return
				case "cat /proc/cpuinfo":
					if _, err := s.Write([]byte(`processor\t: 0\nvendor_id\t: GenuineIntel\ncpu family\t: 6\nmodel\t\t: 142\nmodel name\t: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz\nstepping\t: 11\nmicrocode\t: 0xde\ncpu MHz\t\t: 1992.002\ncache size\t: 8192 KB\nphysical id\t: 0\nsiblings\t: 8\ncore id\t\t: 0\ncpu cores\t: 4\napicid\t\t: 0\ninitial apicid\t: 0\nfpu\t\t: yes\nfpu_exception\t: yes\ncpuid level\t: 22\nwp\t\t: yes\nflags\t\t: fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc art arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc cpuid aperfmperf pni pclmulqdq dtes64 monitor ds_cpl vmx est tm2 ssse3 sdbg fma cx16 xtpr pdcm pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch cpuid_fault epb invpcid_single ssbd ibrs ibpb stibp ibrs_enhanced tpr_shadow vnmi flexpriority ept vpid ept_ad fsgsbase tsc_adjust bmi1 avx2 smep bmi2 erms invpcid mpx rdseed adx smap clflushopt intel_pt xsaveopt xsavec xgetbv1 xsaves dtherm ida arat pln pts hwp hwp_notify hwp_act_window hwp_epp md_clear flush_l1d\nvmx flags\t: vnmi preemption_timer invvpid ept_x_only ept_ad ept_1gb flexpriority tsc_offset vtpr mtf vapic ept vpid unrestricted_guest pml ept_mode_based_exec\nbugs\t\t: spectre_v1 spectre_v2 mds swapgs taa itlb_multihit srbds mmio_stale_data retbleed\nbogomips\t: 3984.00\nclflush size\t: 64\ncache_alignment\t: 64\naddress sizes\t: 39 bits physical, 48 bits virtual\npower management:\n`)); err != nil {
						return
					}
					return
				case "cat /proc/meminfo":
					if _, err := s.Write([]byte(`MemTotal:       16384000 kB\nMemFree:         8192000 kB\nMemAvailable:    12288000 kB\nBuffers:          1024000 kB\nCached:          4096000 kB\nSwapCached:            0 kB\nActive:          6144000 kB\nInactive:        4096000 kB\nActive(anon):    3072000 kB\nInactive(anon):        0 kB\nActive(file):    3072000 kB\nInactive(file):  4096000 kB\nUnevictable:           0 kB\nMlocked:               0 kB\nSwapTotal:             0 kB\nSwapFree:              0 kB\nDirty:                 0 kB\nWriteback:             0 kB\nAnonPages:       3072000 kB\nMapped:          2048000 kB\nShmem:                 0 kB\nKReclaimable:     512000 kB\nSlab:            1024000 kB\nSReclaimable:     512000 kB\nSUnreclaim:       512000 kB\nKernelStack:       81920 kB\nPageTables:       102400 kB\nNFS_Unstable:          0 kB\nBounce:                0 kB\nWritebackTmp:          0 kB\nCommitLimit:     8192000 kB\nCommitted_AS:    6144000 kB\nVmallocTotal:   34359738367 kB\nVmallocUsed:           0 kB\nVmallocChunk:          0 kB\nPercpu:             2048 kB\nHardwareCorrupted:     0 kB\nAnonHugePages:         0 kB\nShmemHugePages:        0 kB\nShmemPmdMapped:        0 kB\nFileHugePages:         0 kB\nFilePmdMapped:         0 kB\nHugePages_Total:       0\nHugePages_Free:        0\nHugePages_Rsvd:        0\nHugePages_Surp:        0\nHugepagesize:       2048 kB\nHugetlb:               0 kB\nDirectMap4k:     4194304 kB\nDirectMap2M:    16777216 kB\nDirectMap1G:           0 kB\n`)); err != nil {
						return
					}
					return
				case "df -B1 /":
					if _, err := s.Write([]byte(`Filesystem     1B-blocks        Used   Available Use% Mounted on\n/dev/sda1     107374182400  21474836480  85899345920  20% /\n`)); err != nil {
						return
					}
					return
				case "ip -json addr":
					if _, err := s.Write([]byte(`[
  {
    "ifindex": 1,
    "ifname": "lo",
    "flags": ["LOOPBACK", "UP", "LOWER_UP"],
    "mtu": 65536,
    "qdisc": "noqueue",
    "operstate": "UNKNOWN",
    "group": "default",
    "txqlen": 1000,
    "link_type": "loopback",
    "address": "00:00:00:00:00:00",
    "broadcast": "00:00:00:00:00:00",
    "addr_info": [
      {
        "family": "inet",
        "local": "127.0.0.1",
        "prefixlen": 8,
        "scope": "host",
        "label": "lo",
        "valid_life_time": 4294967295,
        "preferred_life_time": 4294967295
      }
    ]
  },
  {
    "ifindex": 2,
    "ifname": "eth0",
    "flags": ["BROADCAST", "MULTICAST", "UP", "LOWER_UP"],
    "mtu": 1500,
    "qdisc": "fq_codel",
    "operstate": "UP",
    "group": "default",
    "txqlen": 1000,
    "link_type": "ether",
    "address": "52:54:00:12:34:56",
    "broadcast": "ff:ff:ff:ff:ff:ff",
    "addr_info": [
      {
        "family": "inet",
        "local": "192.168.1.100",
        "prefixlen": 24,
        "scope": "global",
        "label": "eth0",
        "valid_life_time": 4294967295,
        "preferred_life_time": 4294967295
      }
    ]
  }
]`)); err != nil {
						return
					}
					return
				case "ip -json link":
					if _, err := s.Write([]byte(`[
  {
    "ifindex": 1,
    "ifname": "lo",
    "flags": ["LOOPBACK", "UP", "LOWER_UP"],
    "mtu": 65536,
    "qdisc": "noqueue",
    "operstate": "UNKNOWN",
    "group": "default",
    "txqlen": 1000,
    "link_type": "loopback",
    "address": "00:00:00:00:00:00",
    "broadcast": "00:00:00:00:00:00"
  },
  {
    "ifindex": 2,
    "ifname": "eth0",
    "flags": ["BROADCAST", "MULTICAST", "UP", "LOWER_UP"],
    "mtu": 1500,
    "qdisc": "fq_codel",
    "operstate": "UP",
    "group": "default",
    "txqlen": 1000,
    "link_type": "ether",
    "address": "52:54:00:12:34:56",
    "broadcast": "ff:ff:ff:ff:ff:ff"
  }
]`)); err != nil {
						return
					}
					return
				case "cat /etc/resolv.conf":
					if _, err := s.Write([]byte(`# Generated by NetworkManager
search example.com
nameserver 8.8.8.8
nameserver 8.8.4.4
`)); err != nil {
						return
					}
					return
				}
			}
			// Handle single-argument commands
			switch cmd[0] {
			case "hostname":
				if _, err := s.Write([]byte("test-host\n")); err != nil {
					return
				}
			case "nproc":
				if _, err := s.Write([]byte("8\n")); err != nil {
					return
				}
			case "env":
				if _, err := s.Write([]byte("PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin\nHOME=/home/testuser\nUSER=testuser\nSHELL=/bin/bash\nTERM=xterm-256color\nLANG=en_US.UTF-8\n")); err != nil {
					return
				}
			default:
				if _, err := s.Write([]byte("command not found\n")); err != nil {
					return
				}
			}
		},
	}

	return &MockSSHServer{
		server: s,
		port:   port,
		addr:   fmt.Sprintf("localhost:%d", port),
	}, nil
}

// Start starts the mock SSH server
func (m *MockSSHServer) Start() error {
	go func() {
		if err := m.server.ListenAndServe(); err != nil {
			fmt.Printf("SSH server error: %v\n", err)
		}
	}()

	// Wait a bit for the server to start
	time.Sleep(100 * time.Millisecond)
	return nil
}

// Stop stops the mock SSH server
func (m *MockSSHServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.server.Shutdown(ctx)
}

// GetAddr returns the server address
func (m *MockSSHServer) GetAddr() string {
	return m.addr
}

func TestSSHCollector(t *testing.T) {
	// Create mock SSH server
	mockServer, err := NewMockSSHServer()
	if err != nil {
		t.Fatalf("Failed to create mock SSH server: %v", err)
	}

	// Start the server
	if err := mockServer.Start(); err != nil {
		t.Fatalf("Failed to start mock SSH server: %v", err)
	}
	defer func() {
		if err := mockServer.Stop(); err != nil {
			t.Logf("Failed to stop mock server: %v", err)
		}
	}()

	// Create SSH client
	machine := &config.Machine{
		Name:     "mock",
		Host:     "localhost",
		Port:     mockServer.port,
		User:     "testuser",
		Password: "testpass",
	}
	sshClient, err := ssh.NewSSHClient(machine, 5)
	if err != nil {
		t.Fatalf("Failed to create SSH client: %v", err)
	}
	defer sshClient.Close()

	// Create SSH collector
	collector := NewSSHCollector(sshClient)

	// Test collecting all facts
	t.Run("CollectAllFacts", func(t *testing.T) {
		collection, err := collector.Collect("test-server")
		if err != nil {
			t.Fatalf("Failed to collect facts: %v", err)
		}

		if collection.Server != "test-server" {
			t.Errorf("Expected server 'test-server', got '%s'", collection.Server)
		}

		if len(collection.Facts) == 0 {
			t.Error("Expected to collect some facts")
		}

		// Verify some key facts were collected
		expectedFacts := []string{FactHostname, FactOSName, FactCPUCores}
		for _, factKey := range expectedFacts {
			if _, exists := collection.Facts[factKey]; !exists {
				t.Errorf("Expected fact '%s' to be collected", factKey)
			}
		}
	})

	// Test collecting specific facts
	t.Run("CollectSpecificFacts", func(t *testing.T) {
		specificKeys := []string{FactHostname, FactOSName}
		collection, err := collector.CollectSpecific("test-server", specificKeys)
		if err != nil {
			t.Fatalf("Failed to collect specific facts: %v", err)
		}

		if len(collection.Facts) != len(specificKeys) {
			t.Errorf("Expected %d facts, got %d", len(specificKeys), len(collection.Facts))
		}

		for _, key := range specificKeys {
			if _, exists := collection.Facts[key]; !exists {
				t.Errorf("Expected fact '%s' to be collected", key)
			}
		}
	})

	// Test getting a single fact
	t.Run("GetSingleFact", func(t *testing.T) {
		fact, err := collector.GetFact("test-server", FactHostname)
		if err != nil {
			t.Fatalf("Failed to get fact: %v", err)
		}

		if fact.Key != FactHostname {
			t.Errorf("Expected fact key '%s', got '%s'", FactHostname, fact.Key)
		}

		if fact.Server != "test-server" {
			t.Errorf("Expected server 'test-server', got '%s'", fact.Server)
		}

		if fact.Source != string(SourceSSH) {
			t.Errorf("Expected source '%s', got '%s'", SourceSSH, fact.Source)
		}
	})

	// Test system facts collection
	t.Run("CollectSystemFacts", func(t *testing.T) {
		collection := &FactCollection{
			Server:    "test-server",
			Timestamp: time.Now(),
			Facts:     make(map[string]*Fact),
		}

		err := collector.collectSystemFacts(collection)
		if err != nil {
			t.Fatalf("Failed to collect system facts: %v", err)
		}

		expectedSystemFacts := []string{FactMachineID, FactHostname, FactFQDN}
		for _, factKey := range expectedSystemFacts {
			if _, exists := collection.Facts[factKey]; !exists {
				t.Errorf("Expected system fact '%s' to be collected", factKey)
			}
		}
	})

	// Test OS facts collection
	t.Run("CollectOSFacts", func(t *testing.T) {
		collection := &FactCollection{
			Server:    "test-server",
			Timestamp: time.Now(),
			Facts:     make(map[string]*Fact),
		}

		err := collector.collectOSFacts(collection)
		if err != nil {
			t.Fatalf("Failed to collect OS facts: %v", err)
		}

		expectedOSFacts := []string{FactOSName, FactOSVersion, FactOSDistro, FactOSArch, FactOSKernel}
		for _, factKey := range expectedOSFacts {
			if _, exists := collection.Facts[factKey]; !exists {
				t.Errorf("Expected OS fact '%s' to be collected", factKey)
			}
		}
	})

	// Test hardware facts collection
	t.Run("CollectHardwareFacts", func(t *testing.T) {
		collection := &FactCollection{
			Server:    "test-server",
			Timestamp: time.Now(),
			Facts:     make(map[string]*Fact),
		}

		err := collector.collectHardwareFacts(collection)
		if err != nil {
			t.Fatalf("Failed to collect hardware facts: %v", err)
		}

		expectedHardwareFacts := []string{FactCPUCores, FactCPUModel, FactMemoryTotal, FactMemoryUsed, FactMemoryAvail, FactDiskTotal, FactDiskUsed, FactDiskAvail}
		for _, factKey := range expectedHardwareFacts {
			if _, exists := collection.Facts[factKey]; !exists {
				t.Errorf("Expected hardware fact '%s' to be collected", factKey)
			}
		}
	})

	// Test network facts collection
	t.Run("CollectNetworkFacts", func(t *testing.T) {
		collection := &FactCollection{
			Server:    "test-server",
			Timestamp: time.Now(),
			Facts:     make(map[string]*Fact),
		}

		err := collector.collectNetworkFacts(collection)
		if err != nil {
			t.Fatalf("Failed to collect network facts: %v", err)
		}

		expectedNetworkFacts := []string{FactNetworkIPs, FactNetworkMACs, FactDNS}
		for _, factKey := range expectedNetworkFacts {
			if _, exists := collection.Facts[factKey]; !exists {
				t.Errorf("Expected network fact '%s' to be collected", factKey)
			}
		}
	})

	// Test environment facts collection
	t.Run("CollectEnvironmentFacts", func(t *testing.T) {
		collection := &FactCollection{
			Server:    "test-server",
			Timestamp: time.Now(),
			Facts:     make(map[string]*Fact),
		}

		err := collector.collectEnvironmentFacts(collection)
		if err != nil {
			t.Fatalf("Failed to collect environment facts: %v", err)
		}

		if _, exists := collection.Facts[FactEnvironment]; !exists {
			t.Error("Expected environment fact to be collected")
		}
	})
}

func TestSSHCollectorWithTestContainer(t *testing.T) {
	// Skip this test if we're not in a containerized environment
	if os.Getenv("SPOOKY_TEST_CONTAINER") != "true" {
		t.Skip("Skipping container test - set SPOOKY_TEST_CONTAINER=true to run")
	}

	// This test would use a real test container
	// For now, we'll just verify the test structure
	t.Log("Container-based SSH collector test would run here")
}

func TestSSHCollectorErrorHandling(t *testing.T) {
	// Create SSH collector with nil client to test error handling
	collector := NewSSHCollector(nil)

	// Test collecting facts with nil client
	_, err := collector.Collect("test-server")
	if err == nil {
		t.Error("Expected error when collecting facts with nil SSH client")
	}

	// Test collecting specific facts with nil client
	_, err = collector.CollectSpecific("test-server", []string{FactHostname})
	if err == nil {
		t.Error("Expected error when collecting specific facts with nil SSH client")
	}

	// Test getting fact with nil client
	_, err = collector.GetFact("test-server", FactHostname)
	if err == nil {
		t.Error("Expected error when getting fact with nil SSH client")
	}
}

func TestSSHCollectorParsing(t *testing.T) {
	collector := &SSHCollector{}

	// Test OS release parsing
	t.Run("ParseOSRelease", func(t *testing.T) {
		osRelease := `NAME="Ubuntu"
VERSION="22.04.3 LTS (Jammy Jellyfish)"
ID=ubuntu
ID_LIKE=debian
PRETTY_NAME="Ubuntu 22.04.3 LTS"
VERSION_ID="22.04"
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
UBUNTU_CODENAME=jammy`

		osInfo := collector.parseOSRelease(osRelease)
		if osInfo.Name != "Ubuntu" {
			t.Errorf("Expected OS name 'Ubuntu', got '%s'", osInfo.Name)
		}
		if osInfo.Version != "22.04.3 LTS (Jammy Jellyfish)" {
			t.Errorf("Expected OS version '22.04.3 LTS (Jammy Jellyfish)', got '%s'", osInfo.Version)
		}
		if osInfo.Distribution != "ubuntu" {
			t.Errorf("Expected OS distribution 'ubuntu', got '%s'", osInfo.Distribution)
		}
	})

	// Test memory info parsing
	t.Run("ParseMemInfo", func(t *testing.T) {
		memInfo := `MemTotal:       16384000 kB
MemFree:         8192000 kB
MemAvailable:    12288000 kB
Buffers:          1024000 kB
Cached:          4096000 kB`

		memory := collector.parseMemInfo(memInfo)
		if memory.Total != 16777216000 {
			t.Errorf("Expected total memory 16777216000, got %d", memory.Total)
		}
		if memory.Available != 12582912000 {
			t.Errorf("Expected available memory 12582912000, got %d", memory.Available)
		}
	})

	// Test disk info parsing
	t.Run("ParseDiskInfo", func(t *testing.T) {
		dfOutput := `Filesystem     1B-blocks        Used   Available Use% Mounted on
/dev/sda1     107374182400  21474836480  85899345920  20% /`

		disk := collector.parseDiskInfo(dfOutput)
		if disk.Total != 107374182400 {
			t.Errorf("Expected total disk 107374182400, got %d", disk.Total)
		}
		if disk.Used != 21474836480 {
			t.Errorf("Expected used disk 21474836480, got %d", disk.Used)
		}
		if disk.Available != 85899345920 {
			t.Errorf("Expected available disk 85899345920, got %d", disk.Available)
		}
	})

	// Test IP addresses parsing
	t.Run("ParseIPAddresses", func(t *testing.T) {
		ipOutput := `[
  {
    "ifindex": 1,
    "ifname": "lo",
    "addr_info": [
      {
        "family": "inet",
        "local": "127.0.0.1",
        "prefixlen": 8
      }
    ]
  },
  {
    "ifindex": 2,
    "ifname": "eth0",
    "addr_info": [
      {
        "family": "inet",
        "local": "192.168.1.100",
        "prefixlen": 24
      }
    ]
  }
]`

		ips := collector.parseIPAddresses(ipOutput)
		expectedIPs := []string{"127.0.0.1/8", "192.168.1.100/24"}
		if len(ips) != len(expectedIPs) {
			t.Errorf("Expected %d IP addresses, got %d", len(expectedIPs), len(ips))
		}
		for i, expectedIP := range expectedIPs {
			if i < len(ips) && ips[i] != expectedIP {
				t.Errorf("Expected IP '%s', got '%s'", expectedIP, ips[i])
			}
		}
	})

	// Test MAC addresses parsing
	t.Run("ParseMACAddresses", func(t *testing.T) {
		macOutput := `[
  {
    "ifindex": 1,
    "ifname": "lo",
    "address": "00:00:00:00:00:00"
  },
  {
    "ifindex": 2,
    "ifname": "eth0",
    "address": "52:54:00:12:34:56"
  }
]`

		macs := collector.parseMACAddresses(macOutput)
		expectedMACs := []string{"00:00:00:00:00:00", "52:54:00:12:34:56"}
		if len(macs) != len(expectedMACs) {
			t.Errorf("Expected %d MAC addresses, got %d", len(expectedMACs), len(macs))
		}
		for i, expectedMAC := range expectedMACs {
			if i < len(macs) && macs[i] != expectedMAC {
				t.Errorf("Expected MAC '%s', got '%s'", expectedMAC, macs[i])
			}
		}
	})

	// Test DNS config parsing
	t.Run("ParseDNSConfig", func(t *testing.T) {
		resolvConf := `# Generated by NetworkManager
search example.com
nameserver 8.8.8.8
nameserver 8.8.4.4`

		dns := collector.parseDNSConfig(resolvConf)
		if len(dns.Nameservers) != 2 {
			t.Errorf("Expected 2 nameservers, got %d", len(dns.Nameservers))
		}
		if dns.Search[0] != "example.com" {
			t.Errorf("Expected search domain 'example.com', got '%s'", dns.Search[0])
		}
	})

	// Test environment parsing
	t.Run("ParseEnvironment", func(t *testing.T) {
		envOutput := `PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
HOME=/home/testuser
USER=testuser
SHELL=/bin/bash`

		env := collector.parseEnvironment(envOutput)
		if env["PATH"] != "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" {
			t.Errorf("Expected PATH to be set correctly")
		}
		if env["HOME"] != "/home/testuser" {
			t.Errorf("Expected HOME to be set correctly")
		}
		if env["USER"] != "testuser" {
			t.Errorf("Expected USER to be set correctly")
		}
	})
}
