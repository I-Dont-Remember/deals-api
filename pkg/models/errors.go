package models

// ItemExistsError to help distinguish when we've seen this
type ItemExistsError struct {
}

func (e *ItemExistsError) Error() string {
	return "Item already exists"
}
