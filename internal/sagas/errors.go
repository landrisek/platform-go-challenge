package sagas

type RetrievableError struct {
}

func (e RetrievableError) Error() string {
	return "Retrievable error happened"
}
