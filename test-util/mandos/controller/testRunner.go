package mandoscontroller

import (
	mj "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/model"
	mjparse "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/parse"
)

// TestExecutor describes a component that can run a VM test.
type TestExecutor interface {
	// ExecuteTest executes the test and checks if it passed. Failure is signaled by returning an error.
	ExecuteTest(*mj.Test) error
}

// TestRunner is a component that can run tests, using a provided executor.
type TestRunner struct {
	Executor TestExecutor
	Parser   mjparse.Parser
}

// NewTestRunner creates new TestRunner instance.
func NewTestRunner(executor TestExecutor, fileResolver mjparse.FileResolver) *TestRunner {
	return &TestRunner{
		Executor: executor,
		Parser: mjparse.Parser{
			FileResolver: fileResolver,
		},
	}
}
