package profiles

// Profile represents a test bundle configuration
type Profile struct {
	Name        string `hcl:"name,label"`
	Description string `hcl:"description,label"`

	Container     *ContainerConfig   `hcl:"container,block"`
	SSH           *SSHConfig         `hcl:"ssh,block"`
	SpookyProject *ProjectConfig     `hcl:"spooky_project,block"`
	Templates     []*TemplatesConfig `hcl:"templates,block"`
	Files         []*FilesConfig     `hcl:"files,block"`
}

// ContainerConfig defines container settings
type ContainerConfig struct {
	BaseImage string `hcl:"base_image,attr"`
	Name      string `hcl:"name,attr"`
	IP        string `hcl:"ip,attr"`
	Port      int    `hcl:"port,attr"`
}

// SSHConfig defines SSH server settings
type SSHConfig struct {
	Config      string       `hcl:"config,attr"`
	Port        int          `hcl:"port,attr"`
	AuthMethods []string     `hcl:"auth_methods,attr"`
	Users       []string     `hcl:"users,attr"`
	Settings    *SSHSettings `hcl:"settings,block"`
}

// SSHSettings defines SSH server configuration
type SSHSettings struct {
	PasswordAuthentication bool `hcl:"password_authentication,attr"`
	PubkeyAuthentication   bool `hcl:"pubkey_authentication,attr"`
	PermitRootLogin        bool `hcl:"permit_root_login,attr"`
	StrictModes            bool `hcl:"strict_modes,attr"`
	MaxAuthTries           int  `hcl:"max_auth_tries,attr"`
}

// ProjectConfig defines Spooky project settings
type ProjectConfig struct {
	Name        string           `hcl:"name,attr"`
	Description string           `hcl:"description,attr"`
	Machines    *MachinesConfig  `hcl:"machines,block"`
	Variables   *VariablesConfig `hcl:"variables,block"`
	Actions     *ActionsConfig   `hcl:"actions,block"`
}

// MachinesConfig defines machine configurations
type MachinesConfig struct {
	Machines []*MachineConfig `hcl:"machine,block"`
}

// MachineConfig defines a single machine
type MachineConfig struct {
	Name           string                `hcl:"name,attr"`
	Hostname       string                `hcl:"hostname,attr"`
	Port           int                   `hcl:"port,attr"`
	User           string                `hcl:"user,attr"`
	Authentication *AuthenticationConfig `hcl:"authentication,block"`
	Tags           []string              `hcl:"tags,attr"`
}

// AuthenticationConfig defines authentication settings
type AuthenticationConfig struct {
	Method    string           `hcl:"method,attr"`
	Password  *PasswordConfig  `hcl:"password,block"`
	PublicKey *PublicKeyConfig `hcl:"public_key,block"`
}

// PasswordConfig defines password authentication
type PasswordConfig struct {
	Value     string `hcl:"value,attr"`
	Encrypted bool   `hcl:"encrypted,attr"`
}

// PublicKeyConfig defines public key authentication
type PublicKeyConfig struct {
	PublicKeyPath string `hcl:"public_key_path,attr"`
}

// VariablesConfig defines project variables
type VariablesConfig struct {
	Variables []*VariableConfig `hcl:"variable,block"`
}

// VariableConfig defines a single variable
type VariableConfig struct {
	Name        string `hcl:"name,attr"`
	Value       string `hcl:"value,attr"`
	Description string `hcl:"description,attr"`
}

// ActionsConfig defines project actions
type ActionsConfig struct {
	Actions []*ActionConfig `hcl:"action,block"`
}

// ActionConfig defines a single action
type ActionConfig struct {
	Name        string   `hcl:"name,attr"`
	Description string   `hcl:"description,attr"`
	Type        string   `hcl:"type,attr"`
	Command     *string  `hcl:"command,attr"`
	Tags        []string `hcl:"tags,attr"`
	// Additional fields for different action types
	Source          *string `hcl:"source,attr"`
	Destination     *string `hcl:"destination,attr"`
	Validate        *bool   `hcl:"validate,attr"`
	SyncSource      *string `hcl:"sync_source,attr"`
	SyncDestination *string `hcl:"sync_destination,attr"`
	SyncMode        *string `hcl:"sync_mode,attr"`
}

// TemplatesConfig defines a templates block
type TemplatesConfig struct {
	Name      string            `hcl:"name,label"`
	Templates []*TemplateConfig `hcl:"template,block"`
}

// TemplateConfig defines a template file
type TemplateConfig struct {
	Name    string `hcl:"name,label"`
	Content string `hcl:"content,attr"`
}

// FilesConfig defines a files block
type FilesConfig struct {
	Name  string        `hcl:"name,label"`
	Files []*FileConfig `hcl:"file,block"`
}

// FileConfig defines a static file
type FileConfig struct {
	Name        string  `hcl:"name,label"`
	Content     string  `hcl:"content,attr"`
	Permissions *string `hcl:"permissions,attr"`
}
