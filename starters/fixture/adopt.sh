#!/usr/bin/env bash
set -euo pipefail

hgmx=$(cd "$(dirname "$0")/../.." && pwd)
dest=${1:?usage: adopt.sh <empty destination directory>}
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
go tool hgmx init apps android \
	--name "Fixture's Notes & More" \
	--origin https://notes.fixture.example \
	--app-id fixture-android \
	--package com.example.fixture.notes
cp -R "$hgmx/starters/fixture/overlay/." .
go tool hgmx skills
go tool hgmx gen
go mod tidy
go tool hgmx credentials
