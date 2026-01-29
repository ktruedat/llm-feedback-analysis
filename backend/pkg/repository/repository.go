package repository

// OptionsWrapper is a generic wrapper for repository options.
// Each microservice can provide their own extension type T for type-safe customization.
type OptionsWrapper[T any] struct {
	Ex  Executor // Executor for repo operations. Used for transactions or specific contexts.
	Ext *T       // repository-specific customization (type-safe)
}

// RepoOption is a function that configures an OptionsWrapper.
type RepoOption[T any] func(*OptionsWrapper[T])

// WithExecutor sets a custom executor for repository operations.
// This allows using a specific transaction or context for the operations.
func WithExecutor[T any](ex Executor) RepoOption[T] {
	return func(ow *OptionsWrapper[T]) {
		ow.Ex = ex
	}
}

// WithExtension sets a repository-specific extension.
func WithExtension[T any](ext *T) RepoOption[T] {
	return func(ow *OptionsWrapper[T]) {
		ow.Ext = ext
	}
}
