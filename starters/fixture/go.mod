module example.com/fixture

go 1.27.0

require github.com/cogikyo/hgmx v0.4.0

//hgmx:ignore hgmx-latest CI builds this fixture against the hgmx checkout
replace github.com/cogikyo/hgmx => ../..

tool github.com/cogikyo/hgmx/cmd/hgmx
