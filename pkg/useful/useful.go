package useful

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CloseError struct {
	err error
}

func NewCloseError(err error) *CloseError {
	return &CloseError{err: err}
}

func (c *CloseError) Error() string {
	return fmt.Sprintf("cant read request.Body: %s", c.err)
}

func (c *CloseError) Unwrap() error {
	return c.err
}

func ReadJSONBody(request *http.Request, dto interface{}) (err error) {
	if request.Header.Get("Content-Type") != "application/json" {
		return errors.New("Content-Type is not application/json")
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return fmt.Errorf("cant read request.Body: %w", err)
	}
	defer func() {
		errdefer := request.Body.Close()
		if errdefer != nil {
			err = errdefer
		}
	}()

	err = json.Unmarshal(body, &dto)
	if err != nil {
		return err
	}
	return nil
}

func WriteJSONBody(response http.ResponseWriter, dto interface{}) (err error) {
	response.Header().Set("Content-Type", "application/json")

	body, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	_, err = response.Write(body)
	if err != nil {
		return err
	}

	return nil
}
