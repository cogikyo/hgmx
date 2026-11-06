set -eu
name=$1 version=$2 commit=$3 url=$4 product=$5
cache="$HOME/Library/Caches/hgmx/tools"
dir="$cache/$name/$version"
if [ -x "$dir/$product" ]; then
	echo "cached"
	exit 0
fi
mkdir -p "$cache"
tmp=$(mktemp -d "$cache/.$name-$version.XXXXXX")
trap 'rm -rf "$tmp"' EXIT
git -c advice.detachedHead=false clone -q --depth 1 --branch "$version" "$url" "$tmp/src" >&2
got=$(git -C "$tmp/src" rev-parse HEAD)
if [ "$got" != "$commit" ]; then
	echo "$name $version resolves to $got, want $commit" >&2
	exit 1
fi
swift build -c release --package-path "$tmp/src" --product "$product" \
	--scratch-path "$tmp/build" --cache-path "$tmp/cache" --config-path "$tmp/config" --security-path "$tmp/security" >&2
mkdir -p "$tmp/out"
cp "$tmp/build/release/$product" "$tmp/out/$product"
mkdir -p "$cache/$name"
mv "$tmp/out" "$dir"
echo "built"
