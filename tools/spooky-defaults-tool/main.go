package main

import (
	"fmt"
	"spooky/internal/schemas"
)

func main() {
	fmt.Println("🔧 Spooky Struct-Based Defaults Tool")
	fmt.Println("=====================================")

	// Create default config generator
	dcg := schemas.NewDefaultConfigGenerator()

	// Get default Spooky config
	fmt.Println("\n📋 Default Spooky Configuration:")
	fmt.Println("--------------------------------")
	spookyDefaults := dcg.GetDefaultSpookyConfig()

	// Show SSH defaults
	fmt.Printf("SSH Timeout: %d seconds\n", spookyDefaults.SSH.Timeout)
	fmt.Printf("SSH Keepalive Interval: %d seconds\n", spookyDefaults.SSH.KeepaliveInterval)
	fmt.Printf("SSH Keepalive Count: %d\n", spookyDefaults.SSH.KeepaliveCount)
	fmt.Printf("SSH Key Scan Timeout: %d seconds\n", spookyDefaults.SSH.KeyScanTimeout)
	fmt.Printf("SSH Known Hosts Strict: %t\n", spookyDefaults.SSH.KnownHostsStrict)
	fmt.Printf("SSH Connection Pool Size: %d\n", spookyDefaults.SSH.ConnectionPoolSize)

	// Show Security defaults
	fmt.Printf("\nSecurity Allow Unsafe Commands: %t\n", spookyDefaults.Security.AllowUnsafeCommands)
	fmt.Printf("Security Audit Logging: %t\n", spookyDefaults.Security.AuditLogging)

	// Show Logging defaults
	fmt.Printf("\nLogging Level: %s\n", spookyDefaults.Logging.Level)
	fmt.Printf("Logging Format: %s\n", spookyDefaults.Logging.Format)
	fmt.Printf("Logging Output: %s\n", spookyDefaults.Logging.Output)
	fmt.Printf("Logging File Permissions: %s\n", spookyDefaults.Logging.FilePerms)
	fmt.Printf("Logging File Append: %t\n", spookyDefaults.Logging.FileAppend)

	// Generate HCL
	fmt.Println("\n📝 Generated HCL Configuration:")
	fmt.Println("--------------------------------")
	hclConfig, err := dcg.ToHCLWithBlockName(spookyDefaults, "spooky")
	if err != nil {
		fmt.Printf("❌ Error generating HCL: %v\n", err)
		return
	}

	fmt.Println(hclConfig)

	// Show other default configs
	fmt.Println("\n🔧 Other Default Configurations:")
	fmt.Println("--------------------------------")

	// Project defaults
	projectDefaults := dcg.GetDefaultProjectConfig()
	fmt.Printf("Project Max Parallel: %d\n", projectDefaults.RunMaxParallel)
	fmt.Printf("Project Validate Before Run: %t\n", projectDefaults.RunValidateBeforeRun)
	fmt.Printf("Project Facts Timeout: %d seconds\n", projectDefaults.FactsTimeout)

	// Logging defaults
	loggingDefaults := dcg.GetDefaultLoggingConfig()
	fmt.Printf("Logging Performance Buffer Size: %d bytes\n", loggingDefaults.PerformanceBufferSize)
	fmt.Printf("Logging Rotation Max Backups: %d\n", loggingDefaults.RotationMaxBackups)

	fmt.Println("\n✅ All defaults extracted from struct tags successfully!")
}
