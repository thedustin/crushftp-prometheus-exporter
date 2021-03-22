# [CrushFTP Prometheus Exporter](https://github.com/thedustin/crushftp-prometheus-exporter)

by [Dustin Breuer](https://github.com/thedustin)

CrushFTP Prometheus Exporter is an exporter for various metrics about
[CrushFTP](https://www.crushftp.com/), written in Go.

The exporter uses an admin account to scrabe the metrics from the
admin web interface (or in other words, it uses the official
["API"](https://www.crushftp.com/crush9wiki/Wiki.jsp?page=API)).

Tested with CrushFTP 9.3.x & 9.4.x.


# Prerequisites

Before installing this project you need:

* 🐀 Go (at least 1.16)


## Installation

To build this project on your local machine, just run `go build`.

It's recommended to set the version and other build informations, e.g. like this:

```bash
go build \
    -ldflags "\
        -X 'github.com/prometheus/common/version.Branch=$(git branch --show-current)' \
        -X 'github.com/prometheus/common/version.BuildDate=$(date)' \
        -X 'github.com/prometheus/common/version.BuildUser=$(whoami)' \
        -X 'github.com/prometheus/common/version.Revision=$(git rev-parse --short HEAD)' \
        -X 'github.com/prometheus/common/version.Version=0.0.1' \
    " \
    -o crushftp-prometheus-exporter .
```

To verify that your build works you can run `./crushftp-prometheus-exporter --version` to print the version info.


## Configuration

Use `./crushftp-prometheus-exporter --help` to print the supported parameters and environment variables.

By default the exporter listens on port 9100 and provides the metrics under the path `/metrics`.

| Option | Default | Environment variable | Description |
| ------ | ------ | ------ | ------ |
| `-h, --help` |  |  | Show context-sensitive help |
| `-l, --listen-address` | `":9100"` | `LISTEN_ADDRESS` | Address to listen on |
| `    --metrics-endpoint` | `"/metrics"` | `METRICS_ENDPOINT` | Path under which to expose metrics |
| `-H, --crushftp-url` | `"http://localhost"` | `CRUSHFTP_URL` | Base URL to the CrushFTP http(s) server |
| `-u, --crushftp-username` |  | `CRUSHFTP_USERNAME` | CrushFTP admin username<br/><br/>⚠️ It's recommended to use a restricted admin user for monitoring. You can configure the admin restrictions for an user in the user managment section "Admin" > "Setup Roles". See [UserManagerAdminRestricted](https://www.crushftp.com/crush9wiki/Wiki.jsp?page=UserManagerAdminRestricted) for more information |
| `-p, --crushftp-password` |  | `CRUSHFTP_PASSWORD` | CrushFTP admin password<br/>It's recommended to pass this as environment variable |
| `-l, --crushftp-insecure` | `false` | `CRUSHFTP_INSECURE` | Ignore server certificate if using https |
| `    --debug` | `false` |  | Enables debug mode and increases logging |
| `-V, --version` |  |  | Display the application version and exit |


## Contributing to this project

If you have suggestions for improving the prometheus exporter, please
[open an issue or pull request on GitHub](https://github.com/thedustin/crushftp-prometheus-exporter/).
