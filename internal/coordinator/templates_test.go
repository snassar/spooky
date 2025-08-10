package coordinator

import (
	"testing"
	"time"

	spookyinterfaces "spooky/internal/interfaces"
	spookylogging "spooky/internal/logging"
	spookyloggingtypes "spooky/internal/types/logging"
	spookytypes "spooky/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestNewCoordinatorTemplatesIntegration(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorTemplatesIntegration(nil, logger)

	assert.NotNil(t, integration)
	assert.NotNil(t, integration.logger)
}

func TestLoadTemplates(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorTemplatesIntegration(nil, logger)

	projectPath := "/test/project"
	context, err := integration.LoadTemplates(projectPath)

	// Should handle nil templates manager gracefully
	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Equal(t, projectPath, context.ProjectPath)
}

func TestValidateTemplate(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorTemplatesIntegration(nil, logger)

	template := &spookytemplatestypes.Template{
		Name:   "test-template",
		Source: "Hello {{name}}!",
	}

	context := &spookyinterfaces.TemplatesContext{
		Templates:     make(map[string]*spookytemplatestypes.Template),
		RenderedCache: make(map[string]string),
		Functions:     make(map[string]interface{}),
	}

	err := integration.ValidateTemplate(template, context)
	assert.NoError(t, err)
}

func TestRenderTemplate(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorTemplatesIntegration(nil, logger)

	template := &spookytemplatestypes.Template{
		Name:   "test-template",
		Source: "Hello {{name}}!",
	}

	context := &spookyinterfaces.TemplatesContext{
		Templates:     make(map[string]*spookytemplatestypes.Template),
		RenderedCache: make(map[string]string),
		Functions:     make(map[string]interface{}),
	}

	result, err := integration.RenderTemplate(template, context)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestCacheTemplate(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorTemplatesIntegration(nil, logger)

	template := &spookytemplatestypes.Template{
		Name:   "test-template",
		Source: "Hello {{name}}!",
	}

	err := integration.CacheTemplate(template)
	assert.NoError(t, err)
}

func TestTemplatesIntegrationPerformance(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorTemplatesIntegration(nil, logger)

	// Test performance requirement: template rendering < 50ms
	template := &spookytemplatestypes.Template{
		Name:   "test-template",
		Source: "Hello {{name}}!",
	}

	context := &spookyinterfaces.TemplatesContext{
		Templates:     make(map[string]*spookytemplatestypes.Template),
		RenderedCache: make(map[string]string),
		Functions:     make(map[string]interface{}),
	}

	start := time.Now()
	_, err := integration.RenderTemplate(template, context)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Less(t, duration, 50*time.Millisecond, "Template rendering should complete within 50ms")
}

func TestTemplatesIntegrationConcurrent(t *testing.T) {
	logger := spookylogging.NewLogger(spookyloggingtypes.Config{Level: spookyloggingtypes.InfoLevel})
	integration := NewCoordinatorTemplatesIntegration(nil, logger)

	// Test concurrent operations
	const numOperations = 10
	results := make(chan error, numOperations)

	for i := 0; i < numOperations; i++ {
		go func() {
			template := &spookytemplatestypes.Template{
				Name:   "test-template",
				Source: "Hello {{name}}!",
			}

			context := &spookyinterfaces.TemplatesContext{
				Templates:     make(map[string]*spookytemplatestypes.Template),
				RenderedCache: make(map[string]string),
				Functions:     make(map[string]interface{}),
			}

			_, err := integration.RenderTemplate(template, context)
			results <- err
		}()
	}

	// Collect results
	for i := 0; i < numOperations; i++ {
		err := <-results
		assert.NoError(t, err)
	}
}
