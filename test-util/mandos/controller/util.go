package mandoscontroller

import mjparse "github.com/kalyan3104/dme-vm-util/test-util/mandos/json/parse"

// NewDefaultFileResolver yields a new DefaultFileResolver instance.
// Reexported here to avoid having all external packages importing the parser.
// DefaultFileResolver is in parse for local tests only.
func NewDefaultFileResolver() *mjparse.DefaultFileResolver {
	return mjparse.NewDefaultFileResolver()
}
