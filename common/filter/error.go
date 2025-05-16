package filter

type FilterError error

type NotFullError struct{}

func (NotFullError) Error() string {
	return "Filter is not full yet, more steps are needed before any elements can leave the filter"
}

type KeyNotFoundError struct{}

func (KeyNotFoundError) Error() string {
	return "Key not found"
}
