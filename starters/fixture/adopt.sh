#!/usr/bin/env bash
set -euo pipefail

hgmx=$(cd "$(dirname "$0")/../.." && pwd)
dest=${1:?usage: adopt.sh <empty destination directory> [android] [ios]}
shift
platforms=(android ios)
[ $# -eq 0 ] || platforms=("$@")
mkdir -p "$dest"
dest=$(cd "$dest" && pwd)
if [ -n "$(ls -A "$dest")" ]; then
	echo "$dest is not empty" >&2
	exit 1
fi

cd "$hgmx/starters/fixture"
cp -R go.mod main.go api proto "$dest/"
cd "$dest"
go mod edit -replace "github.com/cogikyo/hgmx=$hgmx"
go mod tidy -e
for platform in "${platforms[@]}"; do
	case $platform in
	android) identity=(--package com.example.fixture.notes) ;;
	ios) identity=(--bundle-id com.example.fixture.notes) ;;
	*)
		echo "unknown platform $platform" >&2
		exit 1
		;;
	esac
	go tool hgmx init apps "$platform" \
		--name "Fixture's Notes & More" \
		--origin https://notes.fixture.example \
		--app-id "fixture-$platform" \
		"${identity[@]}"
	cp -R "$hgmx/starters/fixture/overlay/apps/$platform/." "apps/$platform/"
done
go tool hgmx skills
go tool hgmx gen
go mod tidy
go tool hgmx credentials
