package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestMetricsEndpoint(t *testing.T) {
	expected := `# HELP crushftp_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which crushftp_exporter was built.
# TYPE crushftp_exporter_build_info gauge
crushftp_exporter_build_info{branch="",goversion="go1.16.2",revision="",version=""} 1
# HELP crushftp_files_downloaded_total Total number of downloaded files.
# TYPE crushftp_files_downloaded_total counter
crushftp_files_downloaded_total 627308
# HELP crushftp_files_uploaded_total Total number of uploaded files.
# TYPE crushftp_files_uploaded_total counter
crushftp_files_uploaded_total 493131
# HELP crushftp_logins_failed_total Total number of failed logins.
# TYPE crushftp_logins_failed_total counter
crushftp_logins_failed_total 44672
# HELP crushftp_logins_successful_total Total number of successful logins.
# TYPE crushftp_logins_successful_total counter
crushftp_logins_successful_total 7.772218e+06
# HELP crushftp_memory_free_bytes Amount of free memory in bytes.
# TYPE crushftp_memory_free_bytes gauge
crushftp_memory_free_bytes 4.86523928e+08
# HELP crushftp_memory_max_bytes Amount of maximum memory in bytes.
# TYPE crushftp_memory_max_bytes gauge
crushftp_memory_max_bytes 6.442450944e+09
# HELP crushftp_memory_used_bytes Amount of used memory in bytes.
# TYPE crushftp_memory_used_bytes gauge
crushftp_memory_used_bytes 5.955927016e+09
# HELP crushftp_network_bytes_received_total Total number of received bytes.
# TYPE crushftp_network_bytes_received_total counter
crushftp_network_bytes_received_total 3.2223131242e+10
# HELP crushftp_network_bytes_sent_total Total number of sent bytes.
# TYPE crushftp_network_bytes_sent_total counter
crushftp_network_bytes_sent_total 2.7885990912e+10
# HELP crushftp_scrapes_total Total number of CrushFTP scrapes.
# TYPE crushftp_scrapes_total counter
crushftp_scrapes_total 1
# HELP crushftp_security_recent_hammering Was the server recently hammered.
# TYPE crushftp_security_recent_hammering gauge
crushftp_security_recent_hammering 0
# HELP crushftp_server_connected_users Number of users connected to this server.
# TYPE crushftp_server_connected_users gauge
crushftp_server_connected_users{server_name="HTTP",server_type="HTTP"} 7
crushftp_server_connected_users{server_name="HTTPS",server_type="HTTPS"} 3
crushftp_server_connected_users{server_name="SFTP Name",server_type="SFTP"} 1
# HELP crushftp_server_connections_total Number of connections to this server.
# TYPE crushftp_server_connections_total counter
crushftp_server_connections_total{server_name="HTTP",server_type="HTTP"} 930133
crushftp_server_connections_total{server_name="HTTPS",server_type="HTTPS"} 47499
crushftp_server_connections_total{server_name="SFTP Name",server_type="SFTP"} 1.4289802e+07
# HELP crushftp_server_info Number of users connected to this server.
# TYPE crushftp_server_info gauge
crushftp_server_info{hostname="example.org",version="9.2.9"} 1
# HELP crushftp_server_running Is the server running.
# TYPE crushftp_server_running gauge
crushftp_server_running{server_name="HTTP",server_type="HTTP"} 1
crushftp_server_running{server_name="HTTPS",server_type="HTTPS"} 1
crushftp_server_running{server_name="SFTP Name",server_type="SFTP"} 1
# HELP crushftp_server_update_available Is an update available.
# TYPE crushftp_server_update_available gauge
crushftp_server_update_available{available_version="9.4.4"} 1
# HELP crushftp_thread_pool_available Number of available threads.
# TYPE crushftp_thread_pool_available gauge
crushftp_thread_pool_available 122
# HELP crushftp_thread_pool_busy Number of busy threads.
# TYPE crushftp_thread_pool_busy gauge
crushftp_thread_pool_busy 26
# HELP crushftp_thread_pool_max Number of max threads.
# TYPE crushftp_thread_pool_max gauge
crushftp_thread_pool_max 200
# HELP crushftp_up Could the server be reached.
# TYPE crushftp_up gauge
crushftp_up 1
# HELP crushftp_uptime_seconds Number of seconds since the server started.
# TYPE crushftp_uptime_seconds gauge
crushftp_uptime_seconds 1
`

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		f, _ := ioutil.ReadFile("testdata/crushftp/get_dashboard_items/valid.xml")

		res.Write(f)
	}))
	defer func() { testServer.Close() }()

	os.Args = []string{
		os.Args[0],
		"-H", testServer.URL,
		"-u", "user",
		"-p", "s3cr3t",
		"--debug",
	}

	go main()
	var resp *http.Response
	var err error

	// works for now :shrug:
	for i := 0; i < 3; i = i + 1 {
		time.Sleep(10 * time.Millisecond)

		if resp, err = http.DefaultClient.Get("http://:9100/metrics"); err == nil {
			break
		}
	}

	if err != nil {
		t.Fatalf("fetching metrics failed: %q", err)
	}

	r := regexp.MustCompile("crushftp_uptime_seconds [0-9e\\-\\+\\.]+")

	body, _ := io.ReadAll(resp.Body)
	actual := r.ReplaceAllString(string(body), "crushftp_uptime_seconds 1")

	if !cmp.Equal(expected, actual) {
		t.Fatalf("metrics response mismatched:\n%s", cmp.Diff(expected, actual))
	}
}
