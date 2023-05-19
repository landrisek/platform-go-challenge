package sagas

type AssetError struct {
	// description is using only common attribute between assets
	// to identify it for external caller if it error out
	// it is up to caller to set description properly 
	Description string `json:"description"`
	Message     string `json:"message"`
}

type ErrorResponse struct {
	Errors []AssetError `json:"errors"`
}

type RetrievableError struct {

}

func (e RetrievableError) Error() string {
	return "Retrievable error happened"
}
