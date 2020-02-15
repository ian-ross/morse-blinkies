package chassis

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// IntParam extracts an integer URL query parameter.
func IntParam(qs url.Values, k string, dst *int) error {
	s := qs.Get(k)
	if s != "" {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*dst = i
	}
	return nil
}

// ListParams retrieves pagination and query parameters.
func ListParams(r *http.Request, q *string, page *int) error {
	qs, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return errors.New("invalid query parameters")
	}

	*page = 1
	if err := IntParam(qs, "page", page); err != nil {
		return err
	}

	if q != nil {
		*q = qs.Get("q")
	}

	return nil
}
