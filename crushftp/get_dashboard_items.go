package crushftp

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
)

type dashboardItemsResponse struct {
	Status string `xml:"response_status"`
	Data   struct {
		FilesDownloaded          int64             `xml:"downloaded_files"`
		FilesUploaded            int64             `xml:"uploaded_files"`
		Hostname                 string            `xml:"hostname"`
		LoginsFailed             int64             `xml:"failed_logins"`
		LoginsSuccessful         int64             `xml:"successful_logins"`
		RamFree                  int64             `xml:"ram_free"`
		RamMax                   int64             `xml:"ram_max"`
		RecentHammering          bool              `xml:"recent_hammering"`
		ThreadPoolAvailable      int64             `xml:"thread_pool_available"`
		ThreadPoolBusy           int64             `xml:"thread_pool_busy"`
		ThreadPoolMax            int64             `xml:"max_threads"`
		TotalServerBytesReceived int64             `xml:"total_server_bytes_received"`
		TotalServerBytesSent     int64             `xml:"total_server_bytes_sent"`
		VersionInfo              version           `xml:"version_info_str"`
		ServerStartTime          unixDateTime      `xml:"server_start_time"`
		Servers                  []dashboardServer `xml:"server_list>server_list_subitem"`
	} `xml:"response_data>result_value"`
}

type dashboardServer struct {
	ConnectedUsers  int64  `xml:"connected_users"`
	ConnectionNumer int64  `xml:"connection_number"`
	ServerType      string `xml:"serverType"`
	Enabled         bool   `xml:"enabled"`
	Running         bool   `xml:"running"`
}

func (c *client) GetDashboardItems() (*dashboardItemsResponse, error) {
	resp, err := c.command("getDashboardItems")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to fetch data: %s", resp.Status)
	}

	v := &dashboardItemsResponse{}
	s, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	xml.Unmarshal(s, v)

	return v, nil
}
