package types

type Target struct {
	URL        string
	LastResult Result
}

type Result struct {
	StatusCode int
	Body       string
}

type Diff struct {
	Target Target
	Old    Result
	New    Result
}
