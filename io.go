package assessment

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var Version string = "(devel)"

func GetReader(name string) io.ReadCloser {
	strategies := []func(string) io.ReadCloser{tryURL, tryFile, tryText}

	for _, fn := range strategies {
		if r := fn(name); r != nil {
			return r
		}
	}

	return nil
}

func tryURL(name string) io.ReadCloser {
	if u, err := url.Parse(name); err == nil {
		if req, err := http.NewRequest(http.MethodGet, u.String(), nil); err == nil {

			dcl := http.DefaultClient

			req.Header.Set("User-Agent", fmt.Sprintf("Mozilla/5.0+ (compatible; Assessment %s; https://github.com/johnweldon/assessment)", Version))

			if resp, err := dcl.Do(req); err == nil {
				if resp.StatusCode == http.StatusOK {
					return resp.Body
				}
			}
		}
	}

	return nil
}

func tryFile(name string) io.ReadCloser {
	if fi, err := os.Stat(name); err == nil && fi.Size() > 0 {
		if fd, err := os.Open(name); err == nil {
			return fd
		}
	}

	return nil
}

func tryText(v string) io.ReadCloser {
	if len(v) > 3 {
		return io.NopCloser(strings.NewReader(v))
	}

	return nil
}
