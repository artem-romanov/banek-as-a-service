package customerrors

import "fmt"

const ()

type HttpNetworkError struct {
	Err error
	Uri string
}

func (e *HttpNetworkError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf(
			"http network error by uri %s, original error: %s",
			e.Uri,
			e.Err.Error(),
		)
	}
	return fmt.Sprintf("http network error by uri: %s", e.Uri)
}

type NotFoundRequestError struct {
	Uri string
	Err error
}

func (e *NotFoundRequestError) Error() string {
	return fmt.Sprintf("data not found by uri: %s", e.Uri)
}

func (e *NotFoundRequestError) Unwrap() error {
	return e.Err
}

type DownloadRequestError struct {
	Err        error
	StatusCode int
	Uri        string
}

func (e *DownloadRequestError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf(
			"failed request. status code: %d, uri: %s, original error: %s",
			e.StatusCode,
			e.Uri,
			e.Err.Error(),
		)
	}

	return fmt.Sprintf(
		"failed request. status code: %d, uri: %s",
		e.StatusCode,
		e.Uri,
	)
}

func (e *DownloadRequestError) Unwrap() error {
	return e.Err
}

type ParseDataError struct {
	Err error
}

func (e *ParseDataError) Error() string {
	return fmt.Sprintf("parsing data error: %s", e.Err.Error())
}
