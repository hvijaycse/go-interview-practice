package main

import (
	"context"
	"fmt"
	"time"
)

// ContextManager defines a simplified interface for basic context operations
type ContextManager interface {
	// Create a cancellable context from a parent context
	CreateCancellableContext(parent context.Context) (context.Context, context.CancelFunc)

	// Create a context with timeout
	CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc)

	// Add a value to context
	AddValue(parent context.Context, key, value interface{}) context.Context

	// Get a value from context
	GetValue(ctx context.Context, key interface{}) (interface{}, bool)

	// Execute a task with context cancellation support
	ExecuteWithContext(ctx context.Context, task func() error) error

	// Wait for context cancellation or completion
	WaitForCompletion(ctx context.Context, duration time.Duration) error
}

// Simple context manager implementation
type simpleContextManager struct{}

// NewContextManager creates a new context manager
func NewContextManager() ContextManager {
	return &simpleContextManager{}
}

// CreateCancellableContext creates a cancellable context
func (cm *simpleContextManager) CreateCancellableContext(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parent)
}

// CreateTimeoutContext creates a context with timeout
func (cm *simpleContextManager) CreateTimeoutContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}

// AddValue adds a key-value pair to the context
func (cm *simpleContextManager) AddValue(parent context.Context, key, value interface{}) context.Context {
	return context.WithValue(parent, key, value)
}

// GetValue retrieves a value from the context
func (cm *simpleContextManager) GetValue(ctx context.Context, key interface{}) (interface{}, bool) {
	val := ctx.Value(key)

	if val == nil {
		return nil, false
	}
	return val, true

}

// ExecuteWithContext executes a task that can be cancelled via context
func (cm *simpleContextManager) ExecuteWithContext(ctx context.Context, task func() error) error {

	resultChan := make(chan error, 1)

	go func() {
		defer close(resultChan)
		resultChan <- task()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case task_err := <-resultChan:
		return task_err
	}

}

// WaitForCompletion waits for a duration or until context is cancelled
func (cm *simpleContextManager) WaitForCompletion(ctx context.Context, duration time.Duration) error {

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
}

// Helper function - simulate work that can be cancelled
func SimulateWork(ctx context.Context, workDuration time.Duration, description string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(workDuration):
		fmt.Println(description)
		return nil
	}
}

// Helper function - process multiple items with context
func ProcessItems(ctx context.Context, items []string) ([]string, error) {

	output := make([]string, 0, len(items))

	for idx, input := range items {

		select {
		case <-ctx.Done():
			return output, ctx.Err()
		default:
		}

		// Simulate item processing time
		processingTime := time.Millisecond * 50
		if err := SimulateWork(ctx, processingTime, fmt.Sprintf("processing item %d", idx)); err != nil {
			return output, err
		}
		processed := fmt.Sprintf("processed_%v", input)
		output = append(output, processed)
	}
	return output, nil
}
