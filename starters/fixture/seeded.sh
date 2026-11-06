#!/usr/bin/env bash
set -uo pipefail

site=${1:?usage: seeded.sh <adopted fixture> [seed...]}
site=$(cd "$site" && pwd)
shift
seeds=("$@")
[ ${#seeds[@]} -gt 0 ] || seeds=(drift managed pins breaking retained baseline lock verification contract)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
failed=0

expect() {
	local name=$1 want=$2
	shift 2
	local out
	if out=$("$@" 2>&1); then
		echo "FAIL $name: passed but should have failed"
		failed=1
	elif ! grep -qF -- "$want" <<<"$out"; then
		echo "FAIL $name: failed without \"$want\""
		tail -20 <<<"$out"
		failed=1
	else
		echo "PASS $name"
	fi
}

fresh() {
	rm -rf "$work/site"
	cp -a "$site" "$work/site"
	cd "$work/site"
}

gradle() {
	(cd apps/android && ./gradlew --no-daemon --quiet :app:assembleDebug)
}

for seed in "${seeds[@]}"; do
	fresh
	case $seed in
	drift)
		echo "// drift" >>apps/android/app/src/main/generated/fixture/api/v1/NoteServiceClient.kt
		expect drift "NoteServiceClient.kt: out of date" go tool hgmx check
		;;
	managed)
		echo "// edit" >>proto/hgmx/auth/v1/auth.proto
		expect managed "auth.proto: differs from hgmx's copy" go tool hgmx check
		;;
	pins)
		sed -i 's/protobuf-javalite:4.36.2=/protobuf-javalite:4.36.1=/' apps/android/app/gradle.lockfile
		expect pins "protobuf-javalite 4.36.1; generated code needs 4.36.2" go tool hgmx check
		;;
	breaking)
		printf '{"base": "HEAD", "releases": [{"platform": "android", "build": 1, "commit": "%s", "majors": [1]}]}\n' "$(git rev-parse HEAD)" >apps/compatibility.json
		sed -i '/string text = 2;/d' proto/fixture/api/v1/notes.proto
		expect breaking "breaking changes against android build 1" go tool hgmx check
		;;
	retained)
		printf '{"base": "HEAD", "releases": [{"platform": "android", "build": 1, "commit": "%s", "majors": [1]}]}\n' "$(git rev-parse HEAD)" >apps/compatibility.json
		expect retained "no retained-client checks yet" go tool hgmx check
		;;
	baseline)
		printf '{"base": "HEAD", "releases": [{"platform": "android", "build": 1, "commit": "%s", "majors": [1]}]}\n' 0123456789abcdef0123456789abcdef01234567 >apps/compatibility.json
		expect baseline "is not in this clone" go tool hgmx check
		;;
	lock)
		sed -i 's/^com.squareup.okhttp3:okhttp:5.5.1=/com.squareup.okhttp3:okhttp:5.5.0=/' apps/android/app/gradle.lockfile
		expect lock "Dependency Locking" gradle
		;;
	verification)
		sed -i 's/1c5b0837a46d13aec4493b97066c896f0f8f26e5c6b3535bec4355bce2a8e1bb/0000000000000000000000000000000000000000000000000000000000000000/' apps/android/gradle/verification-metadata.xml
		expect verification "Dependency verification failed" gradle
		;;
	contract)
		sed -i 's/string text = 2;/string body = 2;/' proto/fixture/api/v1/notes.proto
		go tool hgmx gen >/dev/null
		expect contract "Unresolved reference" gradle
		;;
	*)
		echo "unknown seed $seed"
		failed=1
		;;
	esac
done
exit $failed
