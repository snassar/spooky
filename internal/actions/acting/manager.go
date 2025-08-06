package acting

import (
	"context"
	"fmt"
	"sync"
	"time"

	"spooky/internal/actions/types"
	"spooky/internal/logging"
)

// Manager implements the ActingManager interface
type Manager struct {
	// Sub-components
	executor        ActingExecutor
	sessionManager  ActingSessionManager
	resultProcessor ActingResultProcessor
	progressTracker ActingProgressTracker

	// Configuration
	defaultTimeout    time.Duration
	defaultParallel   bool
	maxConcurrent     int
	connectionTimeout time.Duration
	commandTimeout    time.Duration
	retryAttempts     int
	retryDelay        time.Duration

	// State
	actors   map[string]Actor
	sessions map[string]*types.ActingSession
	logger   logging.Logger
	mu       sync.RWMutex
}

// NewManager creates a new ActingManager
func NewManager(
	executor ActingExecutor,
	sessionManager ActingSessionManager,
	resultProcessor ActingResultProcessor,
	progressTracker ActingProgressTracker,
	logger logging.Logger,
) *Manager {
	return &Manager{
		executor:          executor,
		sessionManager:    sessionManager,
		resultProcessor:   resultProcessor,
		progressTracker:   progressTracker,
		defaultTimeout:    30 * time.Minute,
		defaultParallel:   false,
		maxConcurrent:     10,
		connectionTimeout: 30 * time.Second,
		commandTimeout:    5 * time.Minute,
		retryAttempts:     3,
		retryDelay:        5 * time.Second,
		actors:            make(map[string]Actor),
		sessions:          make(map[string]*types.ActingSession),
		logger:            logger,
	}
}

// ExecuteAction executes an action
func (m *Manager) ExecuteAction(ctx context.Context, action *types.Action, context *types.ActionContext) (*types.ActingSession, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	if context == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	m.logger.Info("Executing action", logging.String("action", action.Name))

	// Create session
	session, err := m.sessionManager.CreateSession(action.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Start session
	if err := m.sessionManager.StartSession(session.SessionID); err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}

	// Create actor
	actor, err := m.CreateActor(action, context)
	if err != nil {
		m.sessionManager.FailSession(session.SessionID, err)
		return nil, fmt.Errorf("failed to create actor: %w", err)
	}

	// Execute action
	result, err := actor.Execute(ctx, context)
	if err != nil {
		m.sessionManager.FailSession(session.SessionID, err)
		return nil, fmt.Errorf("failed to execute action: %w", err)
	}

	// Process result
	if err := m.resultProcessor.ProcessResult(result); err != nil {
		m.logger.Warn("Failed to process result", logging.Error(err))
	}

	// Add result to session
	session.Results = append(session.Results, result)

	// Complete session
	if err := m.sessionManager.CompleteSession(session.SessionID); err != nil {
		m.logger.Warn("Failed to complete session", logging.Error(err))
	}

	// Update session
	if err := m.sessionManager.UpdateSession(session); err != nil {
		m.logger.Warn("Failed to update session", logging.Error(err))
	}

	m.logger.Info("Action execution completed",
		logging.String("action", action.Name),
		logging.String("session", session.SessionID))

	return session, nil
}

// ExecuteActionCollection executes a collection of actions
func (m *Manager) ExecuteActionCollection(ctx context.Context, collection *types.ActionCollection, context *types.ActionContext) (*types.ActingSession, error) {
	if collection == nil {
		return nil, fmt.Errorf("collection cannot be nil")
	}

	if context == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}

	m.logger.Info("Executing action collection", logging.Int("actions_count", len(collection.Actions)))

	// Create session for collection
	session, err := m.sessionManager.CreateSession("collection")
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Start session
	if err := m.sessionManager.StartSession(session.SessionID); err != nil {
		return nil, fmt.Errorf("failed to start session: %w", err)
	}

	// Start progress tracking
	if err := m.progressTracker.StartTracking(session.SessionID, len(collection.Actions)); err != nil {
		m.logger.Warn("Failed to start progress tracking", logging.Error(err))
	}

	// Execute actions
	var results []*types.ActingResult
	var wg sync.WaitGroup
	resultChan := make(chan *types.ActingResult, len(collection.Actions))
	errorChan := make(chan error, len(collection.Actions))

	// Determine execution strategy
	if m.defaultParallel && len(collection.Actions) > 1 {
		// Parallel execution
		semaphore := make(chan struct{}, m.maxConcurrent)

		for _, action := range collection.Actions {
			wg.Add(1)
			go func(action *types.Action) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				result, err := m.executeActionWithRetry(ctx, action, context)
				if err != nil {
					errorChan <- fmt.Errorf("action '%s' failed: %w", action.Name, err)
					return
				}
				resultChan <- result
			}(action)
		}

		wg.Wait()
		close(resultChan)
		close(errorChan)

		// Collect results
		for result := range resultChan {
			results = append(results, result)
		}

		// Check for errors
		for err := range errorChan {
			m.logger.Error("Action execution failed", err)
		}
	} else {
		// Sequential execution
		for i, action := range collection.Actions {
			result, err := m.executeActionWithRetry(ctx, action, context)
			if err != nil {
				m.logger.Error("Action execution failed", err,
					logging.String("action", action.Name))
				continue
			}
			results = append(results, result)

			// Update progress
			if err := m.progressTracker.UpdateProgress(session.SessionID, i+1); err != nil {
				m.logger.Warn("Failed to update progress", logging.Error(err))
			}
		}
	}

	// Process results
	if err := m.resultProcessor.ProcessResults(results); err != nil {
		m.logger.Warn("Failed to process results", logging.Error(err))
	}

	// Aggregate results
	aggregatedSession, err := m.resultProcessor.AggregateResults(results)
	if err != nil {
		m.logger.Warn("Failed to aggregate results", logging.Error(err))
	} else {
		session = aggregatedSession
	}

	// Add results to session
	session.Results = results

	// Complete session
	if err := m.sessionManager.CompleteSession(session.SessionID); err != nil {
		m.logger.Warn("Failed to complete session", logging.Error(err))
	}

	// Update session
	if err := m.sessionManager.UpdateSession(session); err != nil {
		m.logger.Warn("Failed to update session", logging.Error(err))
	}

	// Complete progress tracking
	if err := m.progressTracker.CompleteTracking(session.SessionID); err != nil {
		m.logger.Warn("Failed to complete progress tracking", logging.Error(err))
	}

	m.logger.Info("Action collection execution completed",
		logging.String("session", session.SessionID),
		logging.Int("total_actions", len(collection.Actions)),
		logging.Int("successful_actions", len(results)))

	return session, nil
}

// PrepareAction prepares an action for execution
func (m *Manager) PrepareAction(action *types.Action, context *types.ActionContext) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	if context == nil {
		return fmt.Errorf("context cannot be nil")
	}

	m.logger.Info("Preparing action", logging.String("action", action.Name))

	// Create actor
	actor, err := m.CreateActor(action, context)
	if err != nil {
		return fmt.Errorf("failed to create actor: %w", err)
	}

	// Prepare actor
	if err := actor.Prepare(context); err != nil {
		return fmt.Errorf("failed to prepare actor: %w", err)
	}

	m.logger.Info("Action prepared", logging.String("action", action.Name))
	return nil
}

// CreateActor creates an actor for an action
func (m *Manager) CreateActor(action *types.Action, context *types.ActionContext) (Actor, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if actor already exists
	if actor, exists := m.actors[action.Name]; exists {
		return actor, nil
	}

	// Create new actor
	actor := &actorImpl{
		action:   action,
		context:  context,
		state:    types.ActingStatePending,
		status:   types.ActingStatusPending,
		timeout:  m.defaultTimeout,
		parallel: m.defaultParallel,
		logger:   m.logger,
	}

	// Store actor
	m.actors[action.Name] = actor

	m.logger.Debug("Created actor", logging.String("action", action.Name))
	return actor, nil
}

// GetActor gets an actor for an action
func (m *Manager) GetActor(action *types.Action) (Actor, error) {
	if action == nil {
		return nil, fmt.Errorf("action cannot be nil")
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	actor, exists := m.actors[action.Name]
	if !exists {
		return nil, fmt.Errorf("actor for action '%s' not found", action.Name)
	}

	return actor, nil
}

// GetSession gets a session by ID
func (m *Manager) GetSession(sessionID string) (*types.ActingSession, error) {
	return m.sessionManager.GetSession(sessionID)
}

// ListSessions lists all sessions
func (m *Manager) ListSessions() ([]*types.ActingSession, error) {
	return m.sessionManager.ListSessions()
}

// CancelSession cancels a session
func (m *Manager) CancelSession(sessionID string) error {
	return m.sessionManager.CancelSession(sessionID)
}

// SetDefaultTimeout sets the default timeout
func (m *Manager) SetDefaultTimeout(timeout time.Duration) {
	m.defaultTimeout = timeout
}

// SetDefaultParallel sets the default parallel flag
func (m *Manager) SetDefaultParallel(parallel bool) {
	m.defaultParallel = parallel
}

// SetMaxConcurrent sets the maximum concurrent executions
func (m *Manager) SetMaxConcurrent(maxConcurrent int) {
	if maxConcurrent < 1 {
		maxConcurrent = 1
	}
	m.maxConcurrent = maxConcurrent
}

// executeActionWithRetry executes an action with retry logic
func (m *Manager) executeActionWithRetry(ctx context.Context, action *types.Action, context *types.ActionContext) (*types.ActingResult, error) {
	var lastErr error

	for attempt := 0; attempt <= m.retryAttempts; attempt++ {
		if attempt > 0 {
			m.logger.Info("Retrying action execution",
				logging.String("action", action.Name),
				logging.Int("attempt", attempt+1),
				logging.Int("max_attempts", m.retryAttempts+1))

			time.Sleep(m.retryDelay)
		}

		// Create actor
		actor, err := m.CreateActor(action, context)
		if err != nil {
			lastErr = err
			continue
		}

		// Execute action
		result, err := actor.Execute(ctx, context)
		if err != nil {
			lastErr = err
			continue
		}

		return result, nil
	}

	return nil, fmt.Errorf("action '%s' failed after %d attempts: %w", action.Name, m.retryAttempts+1, lastErr)
}
