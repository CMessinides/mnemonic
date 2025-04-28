#!/bin/sh

set -e

build_assets() (
    cd internal/assets/src && \
    npx esbuild "*.js" "*.css" \
        --bundle --minify \
        --outdir=../public/dist \
        --entry-names=[dir]/[name]-[hash] \
        --metafile=../public/dist/meta.json
)

build_mnemonicd() {
    go build ./cmd/mnemonicd
}

build_assets
build_mnemonicd
