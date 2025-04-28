#!/bin/sh

set -e

build_ui() (
    cd internal/ui/src && \
    rm ../public/dist/* && \
    npx esbuild "*.js" "*.css" \
        --bundle --minify \
        --outdir=../public/dist \
        --entry-names=[dir]/[name]
)

build_mnemonicd() {
    go build ./cmd/mnemonicd
}

build_ui
build_mnemonicd
