package profiles

// Profile represents a test bundle configuration
type Profile struct {
	Name        string `hcl:"name,label"`
	Description string `hcl:"description,label"`

	DescriptionText string            `hcl:"description,attr"`
	Containers      *ContainersConfig `hcl:"containers,block"`
	Project         *ProjectConfig    `hcl:"project,block"`
}

// ContainersConfig defines multiple container configurations
type ContainersConfig struct {
	Debian13    *ContainerConfig `hcl:"debian13,block"`
	Fedora42    *ContainerConfig `hcl:"fedora42,block"`
	Arch        *ContainerConfig `hcl:"arch,block"`
	Alpine319   *ContainerConfig `hcl:"alpine319,block"`
	Opensuse156 *ContainerConfig `hcl:"opensuse156,block"`
}

// ContainerConfig defines container settings
type ContainerConfig struct {
	BaseImage string     `hcl:"base_image,attr"`
	StaticIP  string     `hcl:"static_ip,attr"`
	SSHPort   int        `hcl:"ssh_port,attr"`
	Packages  []string   `hcl:"packages,attr"`
	SSHConfig *SSHConfig `hcl:"ssh_config,block"`
}

// SSHConfig defines SSH server settings
type SSHConfig struct {
	Port            int    `hcl:"port,attr"`
	PasswordAuth    bool   `hcl:"password_auth,attr"`
	PubkeyAuth      bool   `hcl:"pubkey_auth,attr"`
	PermitRootLogin string `hcl:"permit_root_login,attr"`
	StrictModes     bool   `hcl:"strict_modes,attr"`
	MaxAuthTries    *int   `hcl:"max_auth_tries,attr"`
	LoginGraceTime  *int   `hcl:"login_grace_time,attr"`
}

// ProjectConfig defines Spooky project settings
type ProjectConfig struct {
	Name        string           `hcl:"name,attr"`
	Description string           `hcl:"description,attr"`
	Machines    *MachinesConfig  `hcl:"machines,block"`
	Variables   *VariablesConfig `hcl:"variables,block"`
	Actions     *ActionsConfig   `hcl:"actions,block"`
	Templates   *TemplatesConfig `hcl:"templates,block"`
	Files       *FilesConfig     `hcl:"files,block"`
}

// MachinesConfig defines machine configurations
type MachinesConfig struct {
	Debian13    *MachineConfig `hcl:"debian13,block"`
	Fedora42    *MachineConfig `hcl:"fedora42,block"`
	Arch        *MachineConfig `hcl:"arch,block"`
	Alpine319   *MachineConfig `hcl:"alpine319,block"`
	Opensuse156 *MachineConfig `hcl:"opensuse156,block"`
}

// MachineConfig defines a single machine
type MachineConfig struct {
	Hostname string       `hcl:"hostname,attr"`
	IP       string       `hcl:"ip,attr"`
	Port     int          `hcl:"port,attr"`
	Auth     *AuthConfig  `hcl:"auth,block"`
	Facts    *FactsConfig `hcl:"facts,block"`
}

// AuthConfig defines authentication settings
type AuthConfig struct {
	Password string  `hcl:"password,attr"`
	Username *string `hcl:"username,attr"`
}

// FactsConfig defines facts gathering settings
type FactsConfig struct {
	Basic     bool `hcl:"basic,attr"`
	Enhanced  bool `hcl:"enhanced,attr"`
	Custom    bool `hcl:"custom,attr"`
	Encrypted bool `hcl:"encrypted,attr"`
}

// VariablesConfig defines project variables
type VariablesConfig struct {
	TestType             *string   `hcl:"test_type,attr"`
	OSList               *[]string `hcl:"os_list,attr"`
	ExpectedFacts        *[]string `hcl:"expected_facts,attr"`
	EnhancedFactsTimeout *string   `hcl:"enhanced_facts_timeout,attr"`
	CustomFactsDir       *string   `hcl:"custom_facts_dir,attr"`
	CustomFactsTimeout   *string   `hcl:"custom_facts_timeout,attr"`
	AuthTimeout          *string   `hcl:"auth_timeout,attr"`
	ConnectionRetries    *int      `hcl:"connection_retries,attr"`
	AppName              *string   `hcl:"app_name,attr"`
	AppVersion           *string   `hcl:"app_version,attr"`
	Environment          *string   `hcl:"environment,attr"`
	MaxConnections       *int      `hcl:"max_connections,attr"`
	LogLevel             *string   `hcl:"log_level,attr"`
	SyncMode             *string   `hcl:"sync_mode,attr"`
	IgnorePatterns       *[]string `hcl:"ignore_patterns,attr"`
	SyncInterval         *string   `hcl:"sync_interval,attr"`
}

// ActionsConfig defines project actions
type ActionsConfig struct {
	GatherFacts              *ActionConfig `hcl:"gather_facts,block"`
	VerifyFacts              *ActionConfig `hcl:"verify_facts,block"`
	CompareFacts             *ActionConfig `hcl:"compare_facts,block"`
	SetupCustomFacts         *ActionConfig `hcl:"setup_custom_facts,block"`
	GatherCustomFacts        *ActionConfig `hcl:"gather_custom_facts,block"`
	VerifyCustomFacts        *ActionConfig `hcl:"verify_custom_facts,block"`
	TestCustomFactsExecution *ActionConfig `hcl:"test_custom_facts_execution,block"`
	TestPasswordConnection   *ActionConfig `hcl:"test_password_connection,block"`
	VerifyPasswordAuth       *ActionConfig `hcl:"verify_password_auth,block"`
	TestInvalidPassword      *ActionConfig `hcl:"test_invalid_password,block"`
	GatherFactsViaPassword   *ActionConfig `hcl:"gather_facts_via_password,block"`
	RenderTemplates          *ActionConfig `hcl:"render_templates,block"`
	VerifyTemplates          *ActionConfig `hcl:"verify_templates,block"`
	TestTemplateVariables    *ActionConfig `hcl:"test_template_variables,block"`
	DeployTemplates          *ActionConfig `hcl:"deploy_templates,block"`
	RestartServices          *ActionConfig `hcl:"restart_services,block"`
	SetupSyncDirectories     *ActionConfig `hcl:"setup_sync_directories,block"`
	CreateTestFiles          *ActionConfig `hcl:"create_test_files,block"`
	StartSync                *ActionConfig `hcl:"start_sync,block"`
	VerifySync               *ActionConfig `hcl:"verify_sync,block"`
	TestFileModifications    *ActionConfig `hcl:"test_file_modifications,block"`
	TestConflictResolution   *ActionConfig `hcl:"test_conflict_resolution,block"`
	StopSync                 *ActionConfig `hcl:"stop_sync,block"`
	CleanupSync              *ActionConfig `hcl:"cleanup_sync,block"`
}

// ActionConfig defines a single action
type ActionConfig struct {
	Name        string    `hcl:"name,attr"`
	Description string    `hcl:"description,attr"`
	Command     string    `hcl:"command,attr"`
	Tags        *[]string `hcl:"tags,attr"`
	Parallel    *bool     `hcl:"parallel,attr"`
	Timeout     string    `hcl:"timeout,attr"`
	DependsOn   *[]string `hcl:"depends_on,attr"`
}

// TemplatesConfig defines templates configuration
type TemplatesConfig struct {
	NginxConfig    *TemplateConfig `hcl:"nginx_config,block"`
	AppConfig      *TemplateConfig `hcl:"app_config,block"`
	SystemdService *TemplateConfig `hcl:"systemd_service,block"`
}

// TemplateConfig defines a template file
type TemplateConfig struct {
	Name        string `hcl:"name,attr"`
	Path        string `hcl:"path,attr"`
	Content     string `hcl:"content,attr"`
	Permissions string `hcl:"permissions,attr"`
	Owner       string `hcl:"owner,attr"`
	Group       string `hcl:"group,attr"`
}

// FilesConfig defines files configuration
type FilesConfig struct {
	CustomFactsScript *FileConfig `hcl:"custom_facts_script,block"`
	PythonFactsScript *FileConfig `hcl:"python_facts_script,block"`
	TestData          *FileConfig `hcl:"test_data,block"`
	SyncConfig        *FileConfig `hcl:"sync_config,block"`
}

// FileConfig defines a static file
type FileConfig struct {
	Name        string `hcl:"name,attr"`
	Path        string `hcl:"path,attr"`
	Content     string `hcl:"content,attr"`
	Permissions string `hcl:"permissions,attr"`
	Owner       string `hcl:"owner,attr"`
	Group       string `hcl:"group,attr"`
}
