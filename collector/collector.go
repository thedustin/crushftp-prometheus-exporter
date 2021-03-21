package collector

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/thedustin/crushftp-prometheus-exporter/crushftp"
)

var (
	Namespace = "crushftp"
)

type CollectorOpts struct {
	Logger, ErrorLogger *log.Logger
}

type collector struct {
	crushftpClient *crushftp.Client
	logger         *log.Logger
	errorLogger    *log.Logger

	up           prometheus.Gauge
	scrapesTotal prometheus.Counter

	bytesReceived *prometheus.Desc
	bytesSent     *prometheus.Desc

	filesDownloaded *prometheus.Desc
	filesUploaded   *prometheus.Desc

	info *prometheus.GaugeVec

	loginsFailed     *prometheus.Desc
	loginsSuccessful *prometheus.Desc

	memoryFree prometheus.Gauge
	memoryMax  prometheus.Gauge
	memoryUsed prometheus.Gauge

	securityRecentHammering prometheus.Gauge

	serversConnectedUsers *prometheus.GaugeVec
	serversConnections    *prometheus.Desc
	serversRunning        *prometheus.GaugeVec

	threadPoolAvailable prometheus.Gauge
	threadPoolBusy      prometheus.Gauge
	threadPoolMax       prometheus.Gauge

	updateAvailable *prometheus.GaugeVec
	uptimeSeconds   prometheus.Gauge
}

func NewCollector(crushftpClient *crushftp.Client, opts CollectorOpts) *collector {
	if opts.Logger == nil {
		opts.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	if opts.ErrorLogger == nil {
		opts.ErrorLogger = log.Default()
	}

	return &collector{
		crushftpClient: crushftpClient,
		logger:         opts.Logger,
		errorLogger:    opts.ErrorLogger,

		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "", "up"),
			Help: "Could the server be reached.",
		}),
		scrapesTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(Namespace, "", "scrapes_total"),
			Help: "Total number of CrushFTP scrapes.",
		}),

		bytesReceived: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "network", "bytes_received_total"),
			"Total number of received bytes.",
			[]string{},
			nil,
		),
		bytesSent: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "network", "bytes_sent_total"),
			"Total number of sent bytes.",
			[]string{},
			nil,
		),

		filesDownloaded: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "files", "downloaded_total"),
			"Total number of downloaded files.",
			[]string{},
			nil,
		),
		filesUploaded: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "files", "uploaded_total"),
			"Total number of uploaded files.",
			[]string{},
			nil,
		),

		info: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "server", "info"),
			Help: "Number of users connected to this server.",
		}, []string{"version", "hostname"}),

		loginsFailed: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "logins", "failed_total"),
			"Total number of failed logins.",
			[]string{},
			nil,
		),
		loginsSuccessful: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "logins", "successful_total"),
			"Total number of successful logins.",
			[]string{},
			nil,
		),

		memoryFree: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "memory", "free_bytes"),
			Help: "Amount of free memory in bytes.",
		}),
		memoryMax: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "memory", "max_bytes"),
			Help: "Amount of maximum memory in bytes.",
		}),
		memoryUsed: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "memory", "used_bytes"),
			Help: "Amount of used memory in bytes.",
		}),

		securityRecentHammering: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "security", "recent_hammering"),
			Help: "Was the server recently hammered.",
		}),

		serversConnectedUsers: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "server", "connected_users"),
			Help: "Number of users connected to this server.",
		}, []string{"server_type", "server_name"}),
		serversConnections: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "server", "connections_total"),
			"Number of connections to this server.",
			[]string{"server_type", "server_name"},
			nil,
		),
		serversRunning: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "server", "running"),
			Help: "Is the server running.",
		}, []string{"server_type", "server_name"}),

		threadPoolAvailable: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "thread_pool", "available"),
			Help: "Number of available threads.",
		}),
		threadPoolBusy: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "thread_pool", "busy"),
			Help: "Number of busy threads.",
		}),
		threadPoolMax: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "thread_pool", "max"),
			Help: "Number of max threads.",
		}),

		uptimeSeconds: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "", "uptime_seconds"),
			Help: "Number of seconds since the server started.",
		}),
		updateAvailable: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "server", "update_available"),
			Help: "Is an update available.",
		}, []string{"available_version"}),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	c.logger.Print("Describing collector")

	ch <- c.up.Desc()
	ch <- c.scrapesTotal.Desc()

	ch <- c.bytesReceived
	ch <- c.bytesSent

	ch <- c.filesDownloaded
	ch <- c.filesUploaded

	c.info.Describe(ch)

	ch <- c.loginsFailed
	ch <- c.loginsSuccessful

	ch <- c.memoryFree.Desc()
	ch <- c.memoryUsed.Desc()
	ch <- c.memoryMax.Desc()

	ch <- c.securityRecentHammering.Desc()

	c.serversConnectedUsers.Describe(ch)
	ch <- c.serversConnections
	c.serversRunning.Describe(ch)

	ch <- c.threadPoolAvailable.Desc()
	ch <- c.threadPoolBusy.Desc()
	ch <- c.threadPoolMax.Desc()

	c.updateAvailable.Describe(ch)
	ch <- c.uptimeSeconds.Desc()
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.logger.Print("Collecting data")

	c.up.Set(1)
	c.scrapesTotal.Inc()

	defer func() {
		ch <- c.up
		ch <- c.scrapesTotal
	}()

	resp, err := c.crushftpClient.GetDashboardItems()
	if err != nil {
		c.up.Set(0)

		c.errorLogger.Printf("failed to fetch and decode crushftp stats: %s", err)
		return
	}

	c.logger.Print("Stats fetched")

	ch <- prometheus.MustNewConstMetric(
		c.bytesReceived,
		prometheus.CounterValue,
		float64(resp.Data.BytesReceivedTotal),
	)
	ch <- prometheus.MustNewConstMetric(
		c.bytesSent,
		prometheus.CounterValue,
		float64(resp.Data.BytesSentTotal),
	)

	ch <- prometheus.MustNewConstMetric(
		c.filesDownloaded,
		prometheus.CounterValue,
		float64(resp.Data.FilesDownloaded),
	)
	ch <- prometheus.MustNewConstMetric(
		c.filesUploaded,
		prometheus.CounterValue,
		float64(resp.Data.FilesUploaded),
	)

	c.info.WithLabelValues(resp.Data.VersionInfo.String(), resp.Data.Hostname).Set(1)

	ch <- prometheus.MustNewConstMetric(
		c.loginsFailed,
		prometheus.CounterValue,
		float64(resp.Data.LoginsFailed),
	)
	ch <- prometheus.MustNewConstMetric(
		c.loginsSuccessful,
		prometheus.CounterValue,
		float64(resp.Data.LoginsSuccessful),
	)

	c.memoryFree.Set(float64(resp.Data.RamFree))
	c.memoryMax.Set(float64(resp.Data.RamMax))
	c.memoryUsed.Set(float64(resp.Data.RamMax - resp.Data.RamFree))

	if resp.Data.RecentHammering {
		c.securityRecentHammering.Set(1)
	} else {
		c.securityRecentHammering.Set(0)
	}

	c.threadPoolAvailable.Set(float64(resp.Data.ThreadPoolAvailable))
	c.threadPoolBusy.Set(float64(resp.Data.ThreadPoolBusy))
	c.threadPoolMax.Set(float64(resp.Data.ThreadPoolMax))

	c.uptimeSeconds.Set(time.Since(resp.Data.StartTime.Time).Seconds())

	if resp.Data.UpdateAvailable {
		c.updateAvailable.WithLabelValues(resp.Data.UpdateAvailableVersion.String()).Set(1)
	} else {
		c.updateAvailable.WithLabelValues(resp.Data.UpdateAvailableVersion.String()).Set(0)
	}

	for _, server := range resp.Data.Servers {
		if !server.Enabled {
			continue
		}

		serverName := server.Name
		if serverName == "" {
			serverName = server.ServerType
		}

		labels := []string{server.ServerType, serverName}

		c.serversConnectedUsers.WithLabelValues(labels...).Set(float64(server.ConnectedUsers))

		ch <- prometheus.MustNewConstMetric(
			c.serversConnections,
			prometheus.CounterValue,
			float64(server.ConnectionNumber),
			labels...,
		)

		if server.Running {
			c.serversRunning.WithLabelValues(labels...).Set(1)
		} else {
			c.serversRunning.WithLabelValues(labels...).Set(0)
		}
	}

	c.uptimeSeconds.Collect(ch)
	c.securityRecentHammering.Collect(ch)
	c.updateAvailable.Collect(ch)

	c.memoryFree.Collect(ch)
	c.memoryMax.Collect(ch)
	c.memoryUsed.Collect(ch)

	c.threadPoolAvailable.Collect(ch)
	c.threadPoolBusy.Collect(ch)
	c.threadPoolMax.Collect(ch)

	c.info.Collect(ch)

	c.serversConnectedUsers.Collect(ch)
	c.serversRunning.Collect(ch)

	c.logger.Print("Data collected")
}
