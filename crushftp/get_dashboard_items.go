package crushftp

import (
	"encoding/xml"
	"fmt"
	"io"
)

type DashboardItemsResponse struct {
	Status string             `xml:"response_status"`
	Data   DashboardItemsData `xml:"response_data>result_value"`
}

type DashboardItemsData struct {
	BytesReceivedTotal     int64             `xml:"total_server_bytes_received"`
	BytesSentTotal         int64             `xml:"total_server_bytes_sent"`
	FilesDownloaded        int64             `xml:"downloaded_files"`
	FilesUploaded          int64             `xml:"uploaded_files"`
	Hostname               string            `xml:"hostname"`
	LoginsFailed           int64             `xml:"failed_logins"`
	LoginsSuccessful       int64             `xml:"successful_logins"`
	RamFree                int64             `xml:"ram_free"`
	RamMax                 int64             `xml:"ram_max"`
	RecentHammering        bool              `xml:"recent_hammering"`
	Servers                []DashboardServer `xml:"server_list>server_list_subitem"`
	StartTime              unixDateTime      `xml:"server_start_time"`
	ThreadPoolAvailable    int64             `xml:"thread_pool_available"`
	ThreadPoolBusy         int64             `xml:"thread_pool_busy"`
	ThreadPoolMax          int64             `xml:"max_threads"`
	UpdateAvailable        bool              `xml:"update_available"`
	UpdateAvailableVersion version           `xml:"update_available_version"`
	VersionInfo            version           `xml:"version_info_str"`
}

type DashboardServer struct {
	ConnectedUsers   int64  `xml:"connected_users"`
	ConnectionNumber int64  `xml:"connection_number"`
	Enabled          bool   `xml:"enabled"`
	Name             string `xml:"server_item_name"`
	Running          bool   `xml:"running"`
	ServerType       string `xml:"serverType"`
}

// GetDashboardItems fetches the data used to display the admin web dashboard
// which contains all relevant metrics
func (c *Client) GetDashboardItems() (*DashboardItemsResponse, error) {
	resp, err := c.command("getDashboardItems")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to fetch data: %s", resp.Status)
	}

	v := &DashboardItemsResponse{}

	s, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(s, v); err != nil {
		return nil, err
	}

	return v, nil

}
