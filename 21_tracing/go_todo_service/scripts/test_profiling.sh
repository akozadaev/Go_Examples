#!/usr/bin/env bash
set -euo pipefail

base_url="${PROFILE_BASE_URL:-http://localhost:8080}"
output_dir="${PROFILE_OUTPUT_DIR:-result}"
seconds="${PROFILE_SECONDS:-15}"

mkdir -p "$output_dir"

curl -fsS "$base_url/debug/pprof/heap" -o "$output_dir/heap.pprof"
curl -fsS "$base_url/debug/pprof/allocs" -o "$output_dir/allocs.pprof"
curl -fsS "$base_url/debug/pprof/goroutine?debug=1" -o "$output_dir/goroutine.txt"
curl -fsS "$base_url/debug/pprof/profile?seconds=$seconds" -o "$output_dir/cpu.pprof"

echo "Profiles saved to $output_dir"
echo "CPU profile duration: ${seconds}s"
