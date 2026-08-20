#!/usr/bin/env bash
set -euo pipefail

base_url="${PROFILE_ADMIN_URL:-http://localhost:9090}"
output_dir="${PROFILE_OUTPUT_DIR:-result}"
seconds="${PROFILE_SECONDS:-15}"
pprof_token="${PPROF_TOKEN:-}"
curl_args=(-fsS)
if [[ -n "$pprof_token" ]]; then
  curl_args+=(-H "Authorization: Bearer $pprof_token")
fi

mkdir -p "$output_dir"

curl "${curl_args[@]}" "$base_url/debug/pprof/heap" -o "$output_dir/heap.pprof"
curl "${curl_args[@]}" "$base_url/debug/pprof/allocs" -o "$output_dir/allocs.pprof"
curl "${curl_args[@]}" "$base_url/debug/pprof/goroutine?debug=1" -o "$output_dir/goroutine.txt"
curl "${curl_args[@]}" "$base_url/debug/pprof/profile?seconds=$seconds" -o "$output_dir/cpu.pprof"

echo "Profiles saved to $output_dir"
echo "CPU profile duration: ${seconds}s"
