#!/usr/bin/env bash
set -euo pipefail

hgmx=$(cd "$(dirname "$0")/../.." && pwd)
dest=${1:?usage: bare.sh <empty destination directory> <android|ios>}
platform=${2:?usage: bare.sh <empty destination directory> <android|ios>}
case $platform in
android) identity=(--package com.example.bare) ;;
ios) identity=(--bundle-id com.example.bare) ;;
*)
	echo "unknown platform $platform" >&2
	exit 1
	;;
esac
mkdir -p "$dest"
dest=$(cd "$dest" && pwd)
if [ -n "$(ls -A "$dest")" ]; then
	echo "$dest is not empty" >&2
	exit 1
fi

cp "$hgmx/starters/fixture/go.mod" "$dest/"
cd "$dest"
go mod edit -module example.com/bare -replace "github.com/cogikyo/hgmx=$hgmx"
mkdir -p cmd/server
printf 'package main\n\nfunc main() {}\n' >cmd/server/main.go
go mod tidy -e
go tool hgmx init apps "$platform" --name Bare --origin https://bare.example --app-id "bare-$platform" "${identity[@]}"
go tool hgmx gen
