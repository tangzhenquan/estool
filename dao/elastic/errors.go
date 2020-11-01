package elastic

import "errors"

var (
	CreateIndexError = errors.New("CreateIndex was not acknowledged. Check that timeout value is correct.")
)
