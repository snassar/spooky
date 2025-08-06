package acting

import (
	"spooky/internal/ssh/types"
)

// ActingEngine defines the interface for action execution
type ActingEngine interface {
	ExecuteAction(connection *types.SSHConnection, action *types.SSHAction) (*types.ActionResult, error)
	ExecuteTemplate(connection *types.SSHConnection, template *types.TemplateAction) (*types.ActionResult, error)
	ExecuteSequential(connection *types.SSHConnection, actions []*types.SSHAction) (*types.ActionResult, error)
	ExecuteParallel(connection *types.SSHConnection, actions []*types.SSHAction) (*types.ActionResult, error)
}
