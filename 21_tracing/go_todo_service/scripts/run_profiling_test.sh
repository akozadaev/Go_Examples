#!/usr/bin/env bash
set -euo pipefail

base_url="${PROFILE_BASE_URL:-http://localhost:8080}"
api_url="${PROFILE_API_URL:-$base_url/api/v1/todos}"
output_dir="${PROFILE_OUTPUT_DIR:-result}"
requests="${PROFILE_REQUESTS:-1000}"
concurrency="${PROFILE_CONCURRENCY:-20}"
seconds="${PROFILE_SECONDS:-15}"

mkdir -p "$output_dir"

seq "$requests" | xargs -P "$concurrency" -I '{}' \
	curl -fsS -o /dev/null -H 'X-User-ID: 1' "$api_url"

curl -fsS "$base_url/debug/pprof/heap?gc=1" -o "$output_dir/heap_load.pprof"
curl -fsS "$base_url/debug/pprof/allocs" -o "$output_dir/allocs_load.pprof"
curl -fsS "$base_url/debug/pprof/goroutine?debug=1" -o "$output_dir/goroutine_load.txt"
curl -fsS "$base_url/debug/pprof/profile?seconds=$seconds" -o "$output_dir/cpu_load.pprof" &
profile_pid=$!

seq "$requests" | xargs -P "$concurrency" -I '{}' \
	curl -fsS -o /dev/null -H 'X-User-ID: 1' "$api_url"

wait "$profile_pid"
echo "Load profiles saved to $output_dir"
