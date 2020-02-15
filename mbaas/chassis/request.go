package chassis

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/c2h5oh/datasize"
	"github.com/pkg/errors"
)

const defaultMaxBodySize = 10 * datasize.MB

// ReadBody reads a request body, limiting to a maximum size.
func ReadBody(r *http.Request, limit64 uint64) ([]byte, error) {
	limit := datasize.ByteSize(limit64)
	allowEmpty := false
	if limit == 1 {
		allowEmpty = true
	}
	if limit == 0 || limit == 1 {
		limit = defaultMaxBodySize
	}
	limiter := &io.LimitedReader{r.Body, int64(limit)}
	data, err := ioutil.ReadAll(limiter)
	defer r.Body.Close()
	if err != nil || (!allowEmpty && len(data) == 0) {
		return nil, errors.New("Invalid request body")
	}
	if limiter.N <= 0 {
		return nil, errors.New(fmt.Sprintf("Request body too large (limit is %s)",
			limit.String()))
	}
	return data, nil
}
