package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Generator handles the generation of spooky configuration files
type Generator struct {
	actionsCount  int
	machinesCount int
	outputPath    string
	projectName   string
}

// NewGenerator creates a new generator instance
func NewGenerator() *Generator {
	return &Generator{
		actionsCount:  actionsCount,
		machinesCount: machinesCount,
		outputPath:    outputPath,
		projectName:   projectName,
	}
}

// GenerateProject generates a complete spooky project structure
func (g *Generator) GenerateProject() error {
	// Create project directory structure
	projectDir := g.getProjectPath()
	if err := g.createProjectStructure(projectDir); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	// Generate project.hcl
	projectHCL, err := g.generateProjectConfig()
	if err != nil {
		return fmt.Errorf("failed to generate project config: %w", err)
	}

	// Generate inventory
	inventoryHCL, err := g.generateInventory()
	if err != nil {
		return fmt.Errorf("failed to generate inventory: %w", err)
	}

	// Generate actions (potentially multiple files)
	actionFiles, err := g.generateActionFiles()
	if err != nil {
		return fmt.Errorf("failed to generate actions: %w", err)
	}

	// Write all files
	if err := g.writeFile(filepath.Join(projectDir, "project.hcl"), projectHCL); err != nil {
		return fmt.Errorf("failed to write project.hcl: %w", err)
	}

	if err := g.writeFile(filepath.Join(projectDir, "inventory.hcl"), inventoryHCL); err != nil {
		return fmt.Errorf("failed to write inventory.hcl: %w", err)
	}

	// Write action files
	for filename, content := range actionFiles {
		if err := g.writeFile(filepath.Join(projectDir, "actions", filename), content); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
	}

	fmt.Printf("Generated spooky project '%s' in %s\n", g.projectName, projectDir)
	fmt.Printf("  - %d machines in inventory.hcl\n", g.machinesCount)
	fmt.Printf("  - %d actions across %d files\n", g.actionsCount, len(actionFiles))
	return nil
}

// GenerateInventoryOnly generates only an inventory file
func (g *Generator) GenerateInventoryOnly() error {
	// Create project directory structure
	projectDir := g.getProjectPath()
	if err := g.createProjectStructure(projectDir); err != nil {
		return fmt.Errorf("failed to create project structure: %w", err)
	}

	// Generate project.hcl
	projectHCL, err := g.generateProjectConfig()
	if err != nil {
		return fmt.Errorf("failed to generate project config: %w", err)
	}

	// Generate inventory
	inventoryHCL, err := g.generateInventory()
	if err != nil {
		return fmt.Errorf("failed to generate inventory: %w", err)
	}

	// Write files
	if err := g.writeFile(filepath.Join(projectDir, "project.hcl"), projectHCL); err != nil {
		return fmt.Errorf("failed to write project.hcl: %w", err)
	}

	if err := g.writeFile(filepath.Join(projectDir, "inventory.hcl"), inventoryHCL); err != nil {
		return fmt.Errorf("failed to write inventory.hcl: %w", err)
	}

	fmt.Printf("Generated spooky project '%s' in %s\n", g.projectName, projectDir)
	fmt.Printf("  - %d machines in inventory.hcl\n", g.machinesCount)
	return nil
}

// generateInventory generates the inventory HCL content
func (g *Generator) generateInventory() (string, error) {
	var machines []string

	// Define machine roles and environments for variety
	roles := []string{"web", "database", "cache", "load-balancer", "monitoring", "app", "api", "worker", "storage", "backup"}
	environments := []string{"production", "staging", "development", "testing"}
	regions := []string{"us-east", "us-west", "eu-west", "eu-central", "ap-southeast"}

	for i := 0; i < g.machinesCount; i++ {
		role := roles[i%len(roles)]
		env := environments[i%len(environments)]
		region := regions[i%len(regions)]
		ip := fmt.Sprintf("10.0.%d.%d", (i/255)+1, (i%255)+1)

		machine := fmt.Sprintf(`  machine "server-%d" {
    host = "%s"
    port = 22
    user = "admin"
    password = "password123"
    tags = {
      role = "%s"
      environment = "%s"
      region = "%s"
      instance = "%d"
    }
  }`, i+1, ip, role, env, region, i+1)

		machines = append(machines, machine)
	}

	inventoryHCL := fmt.Sprintf(`# Generated inventory for testing
# Contains %d machines with diverse roles and environments

inventory {
%s
}
`, g.machinesCount, strings.Join(machines, "\n\n"))

	return inventoryHCL, nil
}

// generateProjectConfig generates the project.hcl file
func (g *Generator) generateProjectConfig() (string, error) {
	projectHCL := fmt.Sprintf(`# Generated spooky project configuration
# Project: %s

project "%s" {
  description = "Generated spooky project for testing"
  inventory_file = "inventory.hcl"
  actions_file = "actions.hcl"
  log_file = "logs/spooky.log"
  facts_db = ".facts.db"
}
`, g.projectName, g.projectName)

	return projectHCL, nil
}

// createProjectStructure creates the spooky project directory structure
func (g *Generator) createProjectStructure(projectDir string) error {
	// Create main project directory
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create subdirectories
	subdirs := []string{
		"actions",
		"files",
		"templates",
		"logs",
		".facts.db",
	}

	for _, subdir := range subdirs {
		path := filepath.Join(projectDir, subdir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create %s: %w", subdir, err)
		}
	}

	// Create .gitignore
	gitignore := `# Spooky project files
.facts.db/
logs/
*.log
.DS_Store
`
	if err := g.writeFile(filepath.Join(projectDir, ".gitignore"), gitignore); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}

// getProjectPath returns the path for the project directory
func (g *Generator) getProjectPath() string {
	if g.outputPath == "" {
		return "./spooky-project"
	}
	return g.outputPath
}

// generateActionFiles generates multiple action files based on action count
func (g *Generator) generateActionFiles() (map[string]string, error) {
	actionFiles := make(map[string]string)

	// Determine how many files to create based on action count
	var numFiles int
	if g.actionsCount <= 50 {
		numFiles = 1
	} else if g.actionsCount <= 200 {
		numFiles = 3
	} else if g.actionsCount <= 500 {
		numFiles = 5
	} else {
		numFiles = 8
	}

	// Calculate actions per file
	actionsPerFile := g.actionsCount / numFiles
	remainingActions := g.actionsCount % numFiles

	// Define action categories for file organization
	categories := []struct {
		name        string
		filename    string
		description string
		templates   []actionTemplate
	}{
		{
			name:        "monitoring",
			filename:    "monitoring.hcl",
			description: "System monitoring and health checks",
			templates:   g.getMonitoringTemplates(),
		},
		{
			name:        "deployment",
			filename:    "deployment.hcl",
			description: "Application deployment and updates",
			templates:   g.getDeploymentTemplates(),
		},
		{
			name:        "security",
			filename:    "security.hcl",
			description: "Security and compliance tasks",
			templates:   g.getSecurityTemplates(),
		},
		{
			name:        "maintenance",
			filename:    "maintenance.hcl",
			description: "System maintenance and cleanup",
			templates:   g.getMaintenanceTemplates(),
		},
		{
			name:        "backup",
			filename:    "backup.hcl",
			description: "Backup and recovery operations",
			templates:   g.getBackupTemplates(),
		},
		{
			name:        "network",
			filename:    "network.hcl",
			description: "Network configuration and monitoring",
			templates:   g.getNetworkTemplates(),
		},
		{
			name:        "database",
			filename:    "database.hcl",
			description: "Database operations and maintenance",
			templates:   g.getDatabaseTemplates(),
		},
		{
			name:        "services",
			filename:    "services.hcl",
			description: "Service management and configuration",
			templates:   g.getServiceTemplates(),
		},
	}

	// Generate actions for each file
	actionIndex := 0
	for i := 0; i < numFiles; i++ {
		category := categories[i%len(categories)]

		// Calculate actions for this file
		fileActions := actionsPerFile
		if i < remainingActions {
			fileActions++
		}

		// Generate actions for this category
		var actions []string
		for j := 0; j < fileActions; j++ {
			template := category.templates[j%len(category.templates)]
			actionName := fmt.Sprintf("%s-%d", template.name, actionIndex+1)

			// Convert tags slice to HCL format
			tagStrings := make([]string, len(template.tags))
			for k, tag := range template.tags {
				tagStrings[k] = fmt.Sprintf(`"%s"`, tag)
			}

			action := fmt.Sprintf(`  action "%s" {
    description = "%s"
    command = "%s"
    tags = [%s]
    timeout = %d
    parallel = %t
  }`, actionName, template.description, template.command, strings.Join(tagStrings, ", "), template.timeout, template.parallel)

			actions = append(actions, action)
			actionIndex++
		}

		// Create HCL content for this file
		content := fmt.Sprintf(`# %s
# %s

actions {
%s
}
`, category.description, fmt.Sprintf("Contains %d %s actions", fileActions, category.name), strings.Join(actions, "\n\n"))

		actionFiles[category.filename] = content
	}

	return actionFiles, nil
}

// actionTemplate represents an action template
type actionTemplate struct {
	name        string
	description string
	command     string
	tags        []string
	timeout     int
	parallel    bool
}

// getMonitoringTemplates returns monitoring action templates
func (g *Generator) getMonitoringTemplates() []actionTemplate {
	return []actionTemplate{
		{"check-cpu-usage", "Check CPU usage percentage", "top -bn1 | grep 'Cpu(s)' | awk '{print $2}' | cut -d'%' -f1", []string{"role=web", "role=database", "role=app"}, 300, true},
		{"check-memory-usage", "Check memory usage", "free -h | grep Mem | awk '{print $3}'", []string{"role=web", "role=database", "role=app"}, 300, true},
		{"check-disk-space", "Check disk space usage", "df -h | grep -E '^/dev/' | awk '$5 > 80 {print $0}'", []string{"role=web", "role=database", "role=storage"}, 300, true},
		{"check-load-average", "Check system load average", "cat /proc/loadavg | awk '{print $1, $2, $3}'", []string{"role=web", "role=database", "role=app"}, 300, true},
		{"check-process-count", "Count running processes", "ps aux | wc -l", []string{"role=web", "role=database", "role=app"}, 300, true},
		{"check-uptime", "Check system uptime", "uptime", []string{"role=web", "role=database", "role=app"}, 300, true},
		{"check-disk-iops", "Check disk I/O statistics", "iostat -x 1 1 | grep -v '^$' | tail -n +3", []string{"role=database", "role=storage"}, 300, true},
		{"check-network-connections", "Check active network connections", "netstat -an | grep ESTABLISHED | wc -l", []string{"role=web", "role=database", "role=api"}, 300, true},
	}
}

// getDeploymentTemplates returns deployment action templates
func (g *Generator) getDeploymentTemplates() []actionTemplate {
	return []actionTemplate{
		{"deploy-frontend", "Deploy frontend application", "cd /var/www/app && git pull origin main && npm install && npm run build", []string{"role=web", "environment=production"}, 1800, false},
		{"deploy-backend", "Deploy backend API", "cd /opt/api && git pull origin main && go build && systemctl restart api", []string{"role=api", "environment=production"}, 1200, false},
		{"update-database-schema", "Update database schema", "cd /opt/app && alembic upgrade head", []string{"role=database", "environment=staging"}, 600, false},
		{"restart-application", "Restart application service", "systemctl restart myapp && systemctl status myapp", []string{"role=app", "role=api"}, 300, false},
		{"deploy-configuration", "Deploy configuration files", "cp /tmp/config/* /etc/myapp/ && systemctl reload myapp", []string{"role=web", "role=app"}, 600, false},
		{"update-ssl-certificates", "Update SSL certificates", "certbot renew --quiet && systemctl reload nginx", []string{"role=web", "role=load-balancer"}, 900, false},
		{"deploy-docker-container", "Deploy Docker container", "docker pull myapp:latest && docker-compose up -d", []string{"role=app", "role=api"}, 1200, false},
		{"rollback-deployment", "Rollback to previous deployment", "cd /var/www/app && git reset --hard HEAD~1 && npm run build", []string{"role=web", "environment=production"}, 1800, false},
	}
}

// getSecurityTemplates returns security action templates
func (g *Generator) getSecurityTemplates() []actionTemplate {
	return []actionTemplate{
		{"audit-user-accounts", "Audit user accounts", "awk -F: '$3 >= 1000 && $3 != 65534 {print $1}' /etc/passwd", []string{"role=web", "role=database", "role=app"}, 300, true},
		{"check-ssl-cert-expiry", "Check SSL certificate expiry", "openssl x509 -in /etc/ssl/certs/nginx.crt -noout -dates", []string{"role=web", "role=load-balancer"}, 300, true},
		{"scan-open-ports", "Scan for open ports", "netstat -tlnp | grep LISTEN", []string{"role=web", "role=database", "role=api"}, 300, true},
		{"check-file-permissions", "Check sensitive file permissions", "find /etc -type f -perm /o+w -ls", []string{"role=web", "role=database", "role=app"}, 300, true},
		{"audit-sudo-access", "Audit sudo access", "grep -r '^[^#].*ALL=' /etc/sudoers.d/ /etc/sudoers", []string{"role=web", "role=database", "role=app"}, 300, true},
		{"check-ssh-keys", "Check SSH key permissions", "ls -la ~/.ssh/", []string{"role=web", "role=database", "role=app"}, 300, true},
		{"update-firewall-rules", "Update firewall rules", "iptables -F && iptables-restore < /etc/iptables/rules.v4", []string{"role=web", "role=database", "role=api"}, 600, false},
		{"verify-ssl-configuration", "Verify SSL configuration", "openssl s_client -connect localhost:443 -servername example.com < /dev/null", []string{"role=web", "role=load-balancer"}, 300, true},
	}
}

// getMaintenanceTemplates returns maintenance action templates
func (g *Generator) getMaintenanceTemplates() []actionTemplate {
	return []actionTemplate{
		{"update-system-packages", "Update system packages", "apt update && apt list --upgradable | grep -v '^WARNING' | wc -l", []string{"role=web", "role=database", "role=app"}, 600, false},
		{"clean-temp-files", "Clean temporary files", "find /tmp -type f -mtime +7 -delete", []string{"role=web", "role=database", "role=app"}, 300, true},
		{"rotate-logs", "Rotate log files", "logrotate -f /etc/logrotate.conf", []string{"role=web", "role=database", "role=app"}, 600, false},
		{"clear-application-cache", "Clear application cache", "redis-cli FLUSHALL && echo 'Cache cleared'", []string{"role=cache", "role=app"}, 300, false},
		{"optimize-database", "Optimize database", "mysqlcheck -u root -p --optimize --all-databases", []string{"role=database"}, 1800, false},
		{"clean-old-backups", "Clean old backup files", "find /backups -name '*.tar.gz' -mtime +30 -delete", []string{"role=backup", "role=storage"}, 600, false},
		{"update-cron-jobs", "Update cron jobs", "crontab -l > /tmp/crontab.tmp && crontab /tmp/crontab.tmp", []string{"role=web", "role=database", "role=app"}, 300, false},
		{"check-disk-fragmentation", "Check disk fragmentation", "fsck -N /dev/sda1", []string{"role=storage"}, 300, true},
	}
}

// getBackupTemplates returns backup action templates
func (g *Generator) getBackupTemplates() []actionTemplate {
	return []actionTemplate{
		{"backup-database", "Backup database", "pg_dump -h localhost -U postgres mydb | gzip > /backups/db_$(date +%Y%m%d_%H%M%S).sql.gz", []string{"role=database", "role=backup"}, 1800, false},
		{"backup-config-files", "Backup configuration files", "tar -czf /backups/config_$(date +%Y%m%d).tar.gz /etc/nginx /etc/ssl", []string{"role=web", "role=backup"}, 900, false},
		{"backup-application-data", "Backup application data", "tar -czf /backups/app_$(date +%Y%m%d).tar.gz /var/www/app/data", []string{"role=app", "role=backup"}, 1200, false},
		{"verify-backup-integrity", "Verify backup integrity", "gzip -t /backups/db_$(date +%Y%m%d).sql.gz && echo 'Backup OK'", []string{"role=database", "role=backup"}, 300, true},
		{"sync-backup-to-remote", "Sync backup to remote storage", "rsync -avz /backups/ backup-server:/backups/", []string{"role=backup", "role=storage"}, 3600, false},
		{"create-system-snapshot", "Create system snapshot", "lvcreate -L 10G -s -n snap_$(date +%Y%m%d) /dev/vg0/root", []string{"role=storage", "role=backup"}, 1800, false},
		{"backup-logs", "Backup log files", "tar -czf /backups/logs_$(date +%Y%m%d).tar.gz /var/log", []string{"role=web", "role=database", "role=backup"}, 900, false},
		{"test-backup-restore", "Test backup restore process", "gunzip -c /backups/db_$(date +%Y%m%d).sql.gz | head -100", []string{"role=database", "role=backup"}, 600, true},
	}
}

// getNetworkTemplates returns network action templates
func (g *Generator) getNetworkTemplates() []actionTemplate {
	return []actionTemplate{
		{"check-network-latency", "Check network latency", "ping -c 5 8.8.8.8 | tail -1 | awk '{print $4}' | cut -d'/' -f2", []string{"role=web", "role=load-balancer"}, 300, true},
		{"test-dns-resolution", "Test DNS resolution", "nslookup api.example.com && dig +short api.example.com", []string{"role=web", "role=api", "role=load-balancer"}, 300, true},
		{"check-bandwidth-usage", "Check bandwidth usage", "iftop -t -s 10 -L 100", []string{"role=load-balancer", "role=web"}, 600, true},
		{"monitor-network-connections", "Monitor network connections", "ss -tuln | grep LISTEN", []string{"role=web", "role=database", "role=api"}, 300, true},
		{"check-firewall-status", "Check firewall status", "ufw status verbose", []string{"role=web", "role=database", "role=api"}, 300, true},
		{"test-ssl-connectivity", "Test SSL connectivity", "openssl s_client -connect api.example.com:443 -servername api.example.com", []string{"role=web", "role=api", "role=load-balancer"}, 300, true},
		{"check-routing-table", "Check routing table", "ip route show", []string{"role=load-balancer", "role=web"}, 300, true},
		{"monitor-network-traffic", "Monitor network traffic", "tcpdump -i eth0 -c 100", []string{"role=load-balancer", "role=web"}, 600, true},
	}
}

// getDatabaseTemplates returns database action templates
func (g *Generator) getDatabaseTemplates() []actionTemplate {
	return []actionTemplate{
		{"optimize-database-queries", "Optimize database queries", "psql -c 'SELECT query, calls, total_time FROM pg_stat_statements ORDER BY total_time DESC LIMIT 10'", []string{"role=database"}, 300, true},
		{"check-database-connections", "Check database connections", "psql -c 'SELECT count(*) FROM pg_stat_activity'", []string{"role=database"}, 300, true},
		{"backup-database-incremental", "Incremental database backup", "pg_dump -h localhost -U postgres mydb --exclude-table=temp_* | gzip > /backups/inc_$(date +%Y%m%d_%H%M%S).sql.gz", []string{"role=database", "role=backup"}, 1200, false},
		{"vacuum-database", "Vacuum database", "psql -c 'VACUUM ANALYZE;'", []string{"role=database"}, 1800, false},
		{"check-database-size", "Check database size", "psql -c 'SELECT pg_size_pretty(pg_database_size(current_database()));'", []string{"role=database"}, 300, true},
		{"monitor-slow-queries", "Monitor slow queries", "psql -c 'SELECT query, mean_time FROM pg_stat_statements ORDER BY mean_time DESC LIMIT 5'", []string{"role=database"}, 300, true},
		{"reindex-database", "Reindex database", "psql -c 'REINDEX DATABASE mydb;'", []string{"role=database"}, 3600, false},
		{"check-database-locks", "Check database locks", "psql -c 'SELECT * FROM pg_locks WHERE NOT granted;'", []string{"role=database"}, 300, true},
	}
}

// getServiceTemplates returns service management action templates
func (g *Generator) getServiceTemplates() []actionTemplate {
	return []actionTemplate{
		{"start-web-service", "Start web service", "systemctl start nginx && systemctl status nginx", []string{"role=web"}, 300, false},
		{"stop-web-service", "Stop web service", "systemctl stop nginx", []string{"role=web"}, 300, false},
		{"restart-database-service", "Restart database service", "systemctl restart postgresql", []string{"role=database"}, 600, false},
		{"enable-application-service", "Enable application service", "systemctl enable myapp && systemctl start myapp", []string{"role=app"}, 300, false},
		{"check-service-status", "Check service status", "systemctl status nginx postgresql myapp", []string{"role=web", "role=database", "role=app"}, 300, true},
		{"reload-configuration", "Reload service configuration", "systemctl reload nginx && systemctl reload myapp", []string{"role=web", "role=app"}, 300, false},
		{"disable-old-service", "Disable old service", "systemctl disable oldservice && systemctl stop oldservice", []string{"role=web", "role=app"}, 300, false},
		{"monitor-service-logs", "Monitor service logs", "journalctl -u nginx -f --lines=50", []string{"role=web"}, 600, true},
	}
}

// writeFile writes content to a file, creating directories if needed
func (g *Generator) writeFile(path, content string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Write file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	return nil
}
