# sentinel

sentinel is a small Go agent that runs nmap scans according to configured pools and writes structured results to an output file. It exposes a minimal HTTP API for health checks and re-configuration.

This README explains how to build, configure and run the agent, the pools JSON format, and how scanning and output work.

## Key concepts

- Pool: a named set of hosts, ports, and an interval. Each pool can be scanned once (interval = 0) or periodically.
- Scanner: runs one goroutine per pool, executes `nmap` and sends parsed XML results into a writer which appends JSON results to a file.
- API: a small HTTP server (port 8080) with `/health` and `/configure` endpoints.

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

- `POOLS_FILEPATH` (required): path to the pools JSON file
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

- Health check

```bash
curl http://localhost:8080/health
```

- Configure (replace pools at runtime)

```bash
curl -X POST -H "Content-Type: application/json" --data @pools.json http://localhost:8080/configure
```

The `/configure` handler in the current code stops the running scanner and returns the parsed pools. (Behavior: it calls `scanner.StopScanning()` and returns the new pools. You can extend it to start a new scanner instance with the new pools.)

## Output

By default, results are appended to `scan.json` (or the file specified in `OUTPUT_FILEPATH`). Each line is a JSON object with this shape:

```json
{
  "pool-name": "example-pool",
  "output": {}
}
```

The `output` structure follows the project's `internal/output` types.

## Graceful shutdown and concurrency notes

- The scanner uses one goroutine per pool to run scans and a single writer goroutine that consumes results from a channel and appends to the output file.
- Cancelling scans and ensuring the writer finishes cleanly requires context cancellation and waiting for goroutines (sync.WaitGroup). The current code includes a `StopScanning()` method on the scanner that cancels workers and waits for them.
- When adding or replacing pools at runtime you should stop the old scanner cleanly, then create a new scanner with the updated pools and call `InitScanning()`.

## Troubleshooting

- If no results are written until the program exits, ensure that the writer goroutine is running and the results channel is being closed when workers finish. Also check that file writes are followed by `f.Sync()`.
- If ports are malformed (extra spaces like " 12100"), trim or validate them when parsing pools.

## Extending the project

Ideas and low-risk improvements you can add:

- Make the `/configure` endpoint replace and restart the scanner automatically (start new scanner after `StopScanning`).
- Add validation and normalization of hosts/ports (trim spaces, validate IPs/ranges).
- Add tests for JSON parsing and the `output.ParseRunOutput` XML parsing.
- Add better shutdown handling for the HTTP server (listen for signals and call scanner.StopScanning()).

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
