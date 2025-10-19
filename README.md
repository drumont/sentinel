# sentinel

sentinel is a small Go agent that runs nmap scans according to configured pools and writes structured results to an output file. It exposes a minimal HTTP API for health checks and re-configuration.

This README explains how to build, configure and run the agent, the pools JSON format, and how scanning and output work.

## Key concepts

- Pool: a named set of hosts, ports, and an interval. Each pool can be scanned once (interval = 0) or periodically.
- Scanner: runs one goroutine per pool, executes `nmap` and sends parsed XML results into a writer which appends JSON results to a file. Has state management (RUNNING/STOPPED).
- Services: a service layer that manages the scanner instance and configuration.
- API: a small HTTP server (port 8080) with `/health`, `/configure`, and `/stop` endpoints.

## Prerequisites

- Go (1.20+ recommended)
- nmap installed and available on PATH
- git (optional)

Verify nmap is available:

```bash
nmap --version
```

## Build

From the repository root:

```bash
# build the agent binary
go build -o bin/agent ./cmd/agent

# or run directly
go run ./cmd/agent
```

## Docker Quick Start

Build and run with Docker:

```bash

# create a pools.json file
cat > pools.json << 'EOF'
[
  {
    "name": "example-pool",
    "hosts": ["host1", "host2"],
    "ports": ["80", "8443", "10000"],
    "interval": 60
  }
]
EOF

# run the container with mounted pools file and output directory
docker run -d \
  --name sentinel \
  -p 8080:8080 \
  -v $(pwd)/pools.json:/app/pools.json:ro \
  -v $(pwd)/output:/app/output \
  -e POOLS_FILEPATH=/app/pools.json \
  -e OUTPUT_FILEPATH=/app/output/scan.json \
  ghcr.io/drumont/sentinel:latest

# check logs
docker logs -f sentinel

# stop and remove
docker stop sentinel && docker rm sentinel
```

## Configuration

The agent reads two environment variables at startup:

- `POOLS_FILEPATH` (optional): path to the pools JSON file. If not provided or file doesn't exist, the agent starts without active scans.
- `OUTPUT_FILEPATH` (optional): path to the output file (default: `scan.json`)

Example (macOS / Linux):

```bash
export POOLS_FILEPATH=$(pwd)/scripts/scan.json
export OUTPUT_FILEPATH=$(pwd)/scan.json
./bin/agent
```

## Pools JSON format

The pools file is a JSON array of pool objects. Each pool has:

- `name` (string)
- `hosts` (array of strings)
- `ports` (array of strings) — e.g. `["80", "443", "10000"]`
- `interval` (integer seconds) — `0` means execute once

Example `pools.json`:

```json
[
  {
    "name": "example-pool",
    "hosts": ["host1", "host2"],
    "ports": ["80", "8443", "10000"],
    "interval": 60
  }
]
```

Notes:
- The code expects the file extension to be `.json`.
- Hosts are passed to `nmap` as separate arguments.
- Ports are joined with commas for the `-p` nmap flag.

## How scanning works

For each pool the agent constructs an nmap command roughly like:

```
nmap -Pn -p 80,8443,10000 -oX -  host1 host2
```

- `-Pn`: skip host discovery
- `-p`: ports list
- `-oX -`: write XML to stdout which the agent parses

The agent parses the XML into Go structs and writes JSON results to the output file. Results are written line-delimited (one JSON object per line) and flushed after each write.

## HTTP API

Start the agent and use these endpoints:

- **Health check**

```bash
curl http://localhost:8080/health
# Returns: {"status": "yes"}
```

- **Configure (replace pools at runtime)**

```bash
curl -X POST -H "Content-Type: application/json" --data @pools.json http://localhost:8080/configure
# Returns: HTTP 202 Accepted
```

The `/configure` endpoint stops any running scanner, creates a new scanner with the provided pools, and starts scanning immediately.

- **Stop scanning**

```bash
curl http://localhost:8080/stop
# Returns: HTTP 202 Accepted
```

The `/stop` endpoint gracefully stops all running scans and puts the scanner in STOPPED state.

## Startup Behavior

- If `POOLS_FILEPATH` is provided and the file exists, the agent loads pools at startup and begins scanning immediately.
- If `POOLS_FILEPATH` is empty or the file doesn't exist, the agent starts with no active scans but is ready to accept configuration via the `/configure` endpoint.
- The agent always starts the HTTP API server on port 8080 regardless of whether pools are configured.

## Output

By default, results are appended to `scan.jsonl` (or the file specified in `OUTPUT_FILEPATH`). Each line is a JSON object with this shape:

```json
{
  "pool-name": "example-pool",
  "output": {}
}
```

The `output` structure follows the project's `internal/output` types.

## Graceful shutdown and concurrency notes

- The scanner maintains state (RUNNING/STOPPED) and uses context cancellation for cooperative shutdown of scan goroutines.
- The scanner uses a single writer goroutine that consumes results from a channel and writes JSON to the output file.
- The `/stop` endpoint and `/configure` endpoint both call `StopScanning()` which cancels the context and waits for all goroutines to finish using `sync.WaitGroup`.
- The services layer (`SentinelServices`) manages the scanner lifecycle and ensures proper cleanup when replacing configurations.

## Troubleshooting

- If no results are written, check that the scanner is in RUNNING state and that pools are properly configured.
- If scans don't start at startup, verify that `POOLS_FILEPATH` points to a valid JSON file and check the logs for parsing errors.
- If ports are malformed (extra spaces like " 12100"), trim or validate them when parsing pools.
- Use the `/health` endpoint to verify the API server is running, and `/stop` to gracefully halt scanning operations.

## Extending the project

Ideas and low-risk improvements you can add:

- Add validation and normalization of hosts/ports (trim spaces, validate IPs/ranges).
- Add tests for JSON parsing and the `output.ParseRunOutput` XML parsing.
- Add better shutdown handling for the HTTP server (listen for signals and call `scanner.StopScanning()`).
- Add a `/status` endpoint to report current scanner state and active pools.
- Add metrics and monitoring capabilities (scan duration, success/failure rates).

## Development

Run the agent locally with the provided example pools file:

```bash
export POOLS_FILEPATH=$(pwd)/scripts/scan.json
go run ./cmd/agent
```

Or build and run the binary:

```bash
go build -o bin/agent ./cmd/agent
POOLS_FILEPATH=$(pwd)/scripts/scan.json OUTPUT_FILEPATH=$(pwd)/scan.json ./bin/agent
```

## License & Contributing

This repository does not include a license file. Add one if you plan to share or reuse the code.

Contributions are welcome — open an issue or submit a PR with a clear description and tests where appropriate.
