package fetcher

	
import (
	"fmt"
	"path"
    "testing"
	"net/http"
	"net/http/httptest"
	"os"
	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

// Sample unit test
func Test_Fetch(t *testing.T) {
	mockData1 := "dummy data"
    svr1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, mockData1)
    }))

	mockData2 := "dummy data"
	svr2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, mockData2)
    }))

	dir, err := os.MkdirTemp("", "test_foler")
	assert.NoError(t, err)
	defer os.RemoveAll(dir) // clean up

	var tests = []struct{
		name string
		urls []string
		expectedFiles []string
		meta bool
	}{
		{
			"Should save HTML",
			[]string{svr1.URL, svr2.URL},
			[]string{
				path.Base(svr1.URL)+".html",
				path.Base(svr2.URL)+".html",
			},
			false,
		},
	}

	for _, tt := range tests {
		fetcher := InitFetcher(dir, &tt.meta)

		for _, url := range tt.urls {
			err = fetcher.Fetch(url)
			assert.NoError(t, err)
		}

		files, err := ioutil.ReadDir(dir)
		assert.NoError(t, err)

		assert.Len(t, files, 2, "expected files count")
		assert.Contains(t, tt.expectedFiles, files[0].Name(), "expected file name")
		assert.Contains(t, tt.expectedFiles, files[1].Name(), "expected file name")

		dat, err := os.ReadFile(dir+"/"+tt.expectedFiles[0])
		assert.NoError(t, err)
		assert.Equal(t, string(dat), mockData1, "expected file content")

		dat, err = os.ReadFile(dir+"/"+tt.expectedFiles[1])
		assert.NoError(t, err)
		assert.Equal(t, string(dat), mockData2, "expected file content")
	}
}