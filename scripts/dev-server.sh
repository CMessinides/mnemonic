#!/bin/sh

set -e

watch_assets() (
    cd internal/assets/src && \
    npx esbuild "*.js" "*.css" \
        --bundle --watch=forever \
        --outdir=../public/dist \
        --entry-names=[dir]/[name] \
        --sourcemap
)

watch_server() {
    air
}

(trap 'kill 0' INT; \
    watch_assets & \
    watch_server & \
    wait)
