package ssh

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"spooky/internal/config"
	"spooky/internal/logging"
)

// TemplateActionExecutor handles template action execution
type TemplateActionExecutor struct{}

// NewTemplateActionExecutor creates a new template action executor
func NewTemplateActionExecutor() *TemplateActionExecutor {
	return &TemplateActionExecutor{}
}

// ExecuteAction executes a template action on target machines
func (tae *TemplateActionExecutor) ExecuteAction(action *config.Action, machines []*config.Machine) error {
	logger := logging.GetLogger()

	if action.Template == nil {
		return fmt.Errorf("template configuration is required for template actions")
	}

	logger.Info("Executing template action",
		logging.String("action_name", action.Name),
		logging.String("action_type", action.Type),
		logging.String("template_source", action.Template.Source),
		logging.String("template_destination", action.Template.Destination),
	)

	switch action.Type {
	case "template_deploy":
		return tae.executeTemplateDeploy(action, machines)
	case "template_evaluate":
		return tae.executeTemplateEvaluate(action, machines)
	case "template_validate":
		return tae.executeTemplateValidate(action, machines)
	case "template_cleanup":
		return tae.executeTemplateCleanup(action, machines)
	default:
		return fmt.Errorf("unsupported template action type: %s", action.Type)
	}
}

// executeTemplateDeploy deploys template files to target servers
func (tae *TemplateActionExecutor) executeTemplateDeploy(action *config.Action, machines []*config.Machine) error {
	logger := logging.GetLogger()

	// Validate template file exists
	if _, err := os.Stat(action.Template.Source); os.IsNotExist(err) {
		return fmt.Errorf("template file does not exist: %s", action.Template.Source)
	}

	// Read template file
	templateContent, err := os.ReadFile(action.Template.Source)
	if err != nil {
		return fmt.Errorf("error reading template file %s: %w", action.Template.Source, err)
	}

	// Validate template syntax before deployment
	if err := tae.validateTemplateSyntax(templateContent); err != nil {
		return fmt.Errorf("template syntax validation failed for %s: %w", action.Template.Source, err)
	}

	logger.Info("Deploying template file",
		logging.String("source", action.Template.Source),
		logging.Int("content_size", len(templateContent)),
		logging.Int("target_machines", len(machines)),
	)

	// Deploy to each target machine
	for _, machine := range machines {
		logger.Info("Deploying template to machine",
			logging.String("machine", machine.Name),
			logging.String("destination", action.Template.Destination),
		)

		// Create SSH client
		sshClient, err := NewSSHClient(machine, 30)
		if err != nil {
			logger.Error("Failed to create SSH client", err,
				logging.String("machine", machine.Name))
			continue
		}

		// Execute operations and close client
		func() {
			defer sshClient.Close()

			// Create destination directory if it doesn't exist
			destDir := filepath.Dir(action.Template.Destination)
			if err := tae.createRemoteDirectory(sshClient, destDir); err != nil {
				logger.Error("Failed to create destination directory", err,
					logging.String("machine", machine.Name),
					logging.String("directory", destDir))
				return
			}

			// Check if file already exists and compare content for idempotency
			fileExists, err := tae.remoteFileExists(sshClient, action.Template.Destination)
			if err != nil {
				logger.Warn("Failed to check if file exists, proceeding with deployment",
					logging.String("machine", machine.Name),
					logging.String("file", action.Template.Destination),
					logging.String("error", err.Error()))
			}

			if fileExists {
				// Check if content is different
				contentChanged, err := tae.hasContentChanged(sshClient, action.Template.Destination, templateContent)
				if err != nil {
					logger.Warn("Failed to compare file content, proceeding with deployment",
						logging.String("machine", machine.Name),
						logging.String("file", action.Template.Destination),
						logging.String("error", err.Error()))
				} else if !contentChanged {
					logger.Info("File content unchanged, skipping deployment",
						logging.String("machine", machine.Name),
						logging.String("file", action.Template.Destination))
					return
				}

				// Create backup if requested
				if action.Template.Backup {
					if err := tae.backupRemoteFile(sshClient, action.Template.Destination); err != nil {
						logger.Error("Failed to create backup, aborting deployment", err,
							logging.String("machine", machine.Name),
							logging.String("file", action.Template.Destination))
						return
					}
					logger.Info("Backup created successfully",
						logging.String("machine", machine.Name),
						logging.String("file", action.Template.Destination))
				}
			}

			// Write template file to remote machine
			if err := tae.writeRemoteFile(sshClient, action.Template.Destination, templateContent); err != nil {
				logger.Error("Failed to write template file", err,
					logging.String("machine", machine.Name),
					logging.String("destination", action.Template.Destination))
				return
			}

			// Validate file was written correctly
			if err := tae.validateRemoteFile(sshClient, action.Template.Destination); err != nil {
				logger.Error("File validation failed after deployment", err,
					logging.String("machine", machine.Name),
					logging.String("file", action.Template.Destination))
				return
			}

			// Set file permissions if specified
			if action.Template.Permissions != "" {
				if err := tae.setRemoteFilePermissions(sshClient, action.Template.Destination, action.Template.Permissions); err != nil {
					logger.Error("Failed to set file permissions", err,
						logging.String("machine", machine.Name),
						logging.String("permissions", action.Template.Permissions))
				}
			}

			// Set file owner/group if specified
			if action.Template.Owner != "" || action.Template.Group != "" {
				if err := tae.setRemoteFileOwnership(sshClient, action.Template.Destination, action.Template.Owner, action.Template.Group); err != nil {
					logger.Error("Failed to set file ownership", err,
						logging.String("machine", machine.Name),
						logging.String("owner", action.Template.Owner),
						logging.String("group", action.Template.Group))
				}
			}

			logger.Info("Successfully deployed template to machine",
				logging.String("machine", machine.Name),
				logging.String("destination", action.Template.Destination),
			)
		}()
	}

	return nil
}

// executeTemplateEvaluate evaluates templates on target servers
func (tae *TemplateActionExecutor) executeTemplateEvaluate(action *config.Action, machines []*config.Machine) error {
	logger := logging.GetLogger()

	// Deploy to each target machine
	for _, machine := range machines {
		logger.Info("Evaluating template on machine",
			logging.String("machine", machine.Name),
			logging.String("source", action.Template.Source),
			logging.String("destination", action.Template.Destination),
		)

		// Create SSH client
		sshClient, err := NewSSHClient(machine, 30)
		if err != nil {
			logger.Error("Failed to create SSH client", err,
				logging.String("machine", machine.Name))
			continue
		}

		// Execute operations and close client
		func() {
			defer sshClient.Close()

			// Backup existing file if requested
			if action.Template.Backup {
				if err := tae.backupRemoteFile(sshClient, action.Template.Destination); err != nil {
					logger.Error("Failed to backup existing file", err,
						logging.String("machine", machine.Name),
						logging.String("file", action.Template.Destination))
					return
				}
			}

			// Evaluate template on remote machine
			evaluatedContent, err := tae.evaluateRemoteTemplate(sshClient, action.Template.Source)
			if err != nil {
				logger.Error("Failed to evaluate template", err,
					logging.String("machine", machine.Name),
					logging.String("template", action.Template.Source))
				return
			}

			// Write evaluated content to destination
			if err := tae.writeRemoteFile(sshClient, action.Template.Destination, evaluatedContent); err != nil {
				logger.Error("Failed to write evaluated template", err,
					logging.String("machine", machine.Name),
					logging.String("destination", action.Template.Destination))
				return
			}

			// Validate result if requested
			if action.Template.Validate {
				if err := tae.validateRemoteFile(sshClient, action.Template.Destination); err != nil {
					logger.Error("Template validation failed", err,
						logging.String("machine", machine.Name),
						logging.String("file", action.Template.Destination))
					return
				}
			}

			logger.Info("Successfully evaluated template on machine",
				logging.String("machine", machine.Name),
				logging.String("destination", action.Template.Destination),
			)
		}()
	}

	return nil
}

// executeTemplateValidate validates templates on target servers
func (tae *TemplateActionExecutor) executeTemplateValidate(action *config.Action, machines []*config.Machine) error {
	return tae.executeTemplateOperation(action, machines, "Validating", "validated", func(sshClient *SSHClient, action *config.Action) error {
		return tae.validateRemoteTemplate(sshClient, action.Template.Source)
	})
}

// executeTemplateCleanup removes template files from target servers
func (tae *TemplateActionExecutor) executeTemplateCleanup(action *config.Action, machines []*config.Machine) error {
	return tae.executeTemplateOperation(action, machines, "Cleaning up", "cleaned up", func(sshClient *SSHClient, action *config.Action) error {
		return tae.removeRemoteFile(sshClient, action.Template.Source)
	})
}

// executeTemplateOperation is a helper function to reduce code duplication
func (tae *TemplateActionExecutor) executeTemplateOperation(
	action *config.Action,
	machines []*config.Machine,
	operationName,
	successVerb string,
	operation func(*SSHClient, *config.Action) error,
) error {
	logger := logging.GetLogger()

	for _, machine := range machines {
		logger.Info(operationName+" template on machine",
			logging.String("machine", machine.Name),
			logging.String("template", action.Template.Source),
		)

		sshClient, err := NewSSHClient(machine, 30)
		if err != nil {
			logger.Error("Failed to create SSH client", err,
				logging.String("machine", machine.Name))
			continue
		}

		// Execute operations and close client
		func() {
			defer sshClient.Close()

			if err := operation(sshClient, action); err != nil {
				logger.Error("Template "+operationName+" failed", err,
					logging.String("machine", machine.Name),
					logging.String("template", action.Template.Source))
				return
			}

			logger.Info("Successfully "+successVerb+" template on machine",
				logging.String("machine", machine.Name),
				logging.String("template", action.Template.Source),
			)
		}()
	}

	return nil
}

// Helper methods for remote operations

func (tae *TemplateActionExecutor) createRemoteDirectory(sshClient *SSHClient, dir string) error {
	cmd := fmt.Sprintf("mkdir -p %s", dir)
	_, err := sshClient.ExecuteCommand(cmd)
	return err
}

func (tae *TemplateActionExecutor) writeRemoteFile(sshClient *SSHClient, path string, content []byte) error {
	// Use echo to write content to file
	escapedContent := strings.ReplaceAll(string(content), "'", "'\"'\"'")
	cmd := fmt.Sprintf("echo '%s' > %s", escapedContent, path)
	_, err := sshClient.ExecuteCommand(cmd)
	return err
}

func (tae *TemplateActionExecutor) setRemoteFilePermissions(sshClient *SSHClient, path, permissions string) error {
	cmd := fmt.Sprintf("chmod %s %s", permissions, path)
	_, err := sshClient.ExecuteCommand(cmd)
	return err
}

func (tae *TemplateActionExecutor) setRemoteFileOwnership(sshClient *SSHClient, path, owner, group string) error {
	var cmd string
	switch {
	case owner != "" && group != "":
		cmd = fmt.Sprintf("chown %s:%s %s", owner, group, path)
	case owner != "":
		cmd = fmt.Sprintf("chown %s %s", owner, path)
	case group != "":
		cmd = fmt.Sprintf("chgrp %s %s", group, path)
	default:
		return nil
	}
	_, err := sshClient.ExecuteCommand(cmd)
	return err
}

func (tae *TemplateActionExecutor) backupRemoteFile(sshClient *SSHClient, path string) error {
	backupPath := path + ".backup"
	cmd := fmt.Sprintf("cp %s %s", path, backupPath)
	_, err := sshClient.ExecuteCommand(cmd)
	return err
}

func (tae *TemplateActionExecutor) removeRemoteFile(sshClient *SSHClient, path string) error {
	cmd := fmt.Sprintf("rm -f %s", path)
	_, err := sshClient.ExecuteCommand(cmd)
	return err
}

func (tae *TemplateActionExecutor) validateRemoteFile(sshClient *SSHClient, path string) error {
	// Basic validation - check if file exists and is readable
	cmd := fmt.Sprintf("test -r %s", path)
	_, err := sshClient.ExecuteCommand(cmd)
	return err
}

// evaluateRemoteTemplate evaluates a template on the remote machine
func (tae *TemplateActionExecutor) evaluateRemoteTemplate(sshClient *SSHClient, templatePath string) ([]byte, error) {
	// Read template content from remote machine
	readCmd := fmt.Sprintf("cat %s", templatePath)
	templateContent, err := sshClient.ExecuteCommand(readCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	// Create server-side template functions
	funcMap := template.FuncMap{
		"machineID": func() string {
			return tae.getRemoteFact(sshClient, "cat /etc/machine-id")
		},
		"osVersion": func() string {
			return tae.getRemoteFact(sshClient, "uname -r")
		},
		"hostname": func() string {
			return tae.getRemoteFact(sshClient, "hostname")
		},
		"ipAddress": func() string {
			return tae.getRemoteFact(sshClient, "hostname -I | awk '{print $1}'")
		},
		"diskSpace": func() string {
			return tae.getRemoteFact(sshClient, "df -h / | tail -1 | awk '{print $4}'")
		},
		"memoryInfo": func() string {
			return tae.getRemoteFact(sshClient, "free -h | grep Mem | awk '{print $2}'")
		},
		"fileExists": func(path string) bool {
			result := tae.getRemoteFact(sshClient, fmt.Sprintf("test -f %s && echo 'true' || echo 'false'", path))
			return result == "true"
		},
		"fileContent": func(path string) string {
			return tae.getRemoteFact(sshClient, fmt.Sprintf("cat %s 2>/dev/null || echo ''", path))
		},
		"fileSize": func(path string) string {
			return tae.getRemoteFact(sshClient, fmt.Sprintf("stat -c%%s %s 2>/dev/null || echo '0'", path))
		},
		"fileOwner": func(path string) string {
			return tae.getRemoteFact(sshClient, fmt.Sprintf("stat -c%%U %s 2>/dev/null || echo ''", path))
		},
	}

	// Parse and execute template
	tmpl, err := template.New("remote_template").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var result bytes.Buffer
	if err := tmpl.Execute(&result, nil); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return result.Bytes(), nil
}

// validateRemoteTemplate validates a template on the remote machine
func (tae *TemplateActionExecutor) validateRemoteTemplate(sshClient *SSHClient, templatePath string) error {
	// Read template content from remote machine
	readCmd := fmt.Sprintf("cat %s", templatePath)
	templateContent, err := sshClient.ExecuteCommand(readCmd)
	if err != nil {
		return fmt.Errorf("failed to read template: %w", err)
	}

	// Basic template validation - check for balanced braces
	openBraces := strings.Count(templateContent, "{{")
	closeBraces := strings.Count(templateContent, "}}")

	if openBraces != closeBraces {
		return fmt.Errorf("unbalanced template braces: %d opening, %d closing", openBraces, closeBraces)
	}

	// Try to parse template (without executing)
	funcMap := template.FuncMap{
		"machineID":   func() string { return "" },
		"osVersion":   func() string { return "" },
		"hostname":    func() string { return "" },
		"ipAddress":   func() string { return "" },
		"diskSpace":   func() string { return "" },
		"memoryInfo":  func() string { return "" },
		"fileExists":  func(_ string) bool { return false },
		"fileContent": func(_ string) string { return "" },
		"fileSize":    func(_ string) string { return "0" },
		"fileOwner":   func(_ string) string { return "" },
	}

	_, err = template.New("validation").Funcs(funcMap).Parse(templateContent)
	return err
}

// getRemoteFact executes a command on the remote machine and returns the result
func (tae *TemplateActionExecutor) getRemoteFact(sshClient *SSHClient, command string) string {
	result, err := sshClient.ExecuteCommand(command)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(result)
}

// remoteFileExists checks if a file exists on the remote machine
func (tae *TemplateActionExecutor) remoteFileExists(sshClient *SSHClient, path string) (bool, error) {
	result, err := sshClient.ExecuteCommand(fmt.Sprintf("test -f %s && echo 'exists' || echo 'not_exists'", path))
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(result) == "exists", nil
}

// hasContentChanged compares the content of a remote file with local content
func (tae *TemplateActionExecutor) hasContentChanged(sshClient *SSHClient, path string, localContent []byte) (bool, error) {
	// Get remote file content
	remoteContent, err := sshClient.ExecuteCommand(fmt.Sprintf("cat %s", path))
	if err != nil {
		return true, err // Assume changed if we can't read remote file
	}

	// Compare content
	return string(localContent) != remoteContent, nil
}

// validateTemplateSyntax validates the syntax of a template file
func (tae *TemplateActionExecutor) validateTemplateSyntax(templateContent []byte) error {
	// Create a minimal template with basic functions for validation
	funcMap := template.FuncMap{
		"machineID":   func() string { return "test" },
		"osVersion":   func() string { return "test" },
		"hostname":    func() string { return "test" },
		"ipAddress":   func() string { return "test" },
		"diskSpace":   func() string { return "test" },
		"memoryInfo":  func() string { return "test" },
		"fileExists":  func(path string) bool { return true },
		"fileContent": func(path string) string { return "test" },
		"fileSize":    func(path string) string { return "test" },
		"fileOwner":   func(path string) string { return "test" },
	}

	// Try to parse the template
	_, err := template.New("validation").Funcs(funcMap).Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("template syntax error: %w", err)
	}

	return nil
}
