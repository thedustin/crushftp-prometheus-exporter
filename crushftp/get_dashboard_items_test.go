package crushftp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGetDashboardItemsUnmarshal(t *testing.T) {
	tests := map[string]*DashboardItemsResponse{
		"valid": {
			Status: "OK",
			Data: DashboardItemsData{
				BytesReceivedTotal: 32223131242,
				BytesSentTotal:     27885990912,
				FilesDownloaded:    627308,
				FilesUploaded:      493131,
				Hostname:           "example.org",
				LoginsFailed:       44672,
				LoginsSuccessful:   7772218,
				RamFree:            486523928,
				RamMax:             6442450944,
				RecentHammering:    false,
				Servers: []DashboardServer{
					{ConnectedUsers: 3, ConnectionNumber: 47499, Enabled: true, Running: true, ServerType: "HTTPS"},
					{ConnectedUsers: 1, ConnectionNumber: 14289802, Enabled: true, Name: "SFTP Name", Running: true, ServerType: "SFTP"},
					{ConnectedUsers: 7, ConnectionNumber: 930133, Enabled: true, Running: true, ServerType: "HTTP"},
				},
				StartTime:              unixDateTime{time.Date(2021, time.January, 02, 12, 14, 16, 0, time.FixedZone("CET", int(time.Hour*1/time.Second)))},
				ThreadPoolAvailable:    122,
				ThreadPoolBusy:         26,
				ThreadPoolMax:          200,
				UpdateAvailable:        true,
				UpdateAvailableVersion: version{Major: 9, Minor: 4, Patch: 4},
				VersionInfo:            version{Major: 9, Minor: 2, Patch: 9},
			},
		},
	}

	for testFile, expected := range tests {
		f, _ := io.ReadAll(getTestDataFile(testFile, t))

		testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Write(f)
		}))
		defer func() { testServer.Close() }()

		url, _ := url.Parse(testServer.URL)
		c := NewClient(ClientOptions{
			HostAndPort: url.Host,
			Scheme:      url.Scheme,
		})
		c.httpClient = testServer.Client()

		actual, err := c.GetDashboardItems()

		if err != nil {
			t.Fatalf("failed with error %q", err)
		}

		if !cmp.Equal(expected, actual) {
			t.Fatalf("response for %q mismatched:\n%s", testFile, cmp.Diff(expected, actual))
		}
	}
}

func getTestDataFile(filename string, t *testing.T) io.Reader {
	f, err := os.Open("../testdata/crushftp/get_dashboard_items/" + filename + ".xml")

	if err != nil {
		t.Fatalf("cannot open test file: %s", err)
	}

	return f
}
