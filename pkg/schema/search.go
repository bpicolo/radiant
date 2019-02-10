package schema

// Search represents a search using a particular query definition
type Search struct {
	Query   *Query
	Context interface{}
	From    int
	Size    int
}
