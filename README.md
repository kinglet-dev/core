<p align="center">
  <a href="https://github.com/kinglet-dev/core/actions/workflows/ci.yml"><img src="https://github.com/kinglet-dev/core/actions/workflows/ci.yml/badge.svg?branch=main" alt="CI status"></a>
  <a href="https://github.com/kinglet-dev/core/actions/workflows/badges.yml"><img src="https://raw.githubusercontent.com/kinglet-dev/core/badges/coverage.svg" alt="Test coverage, measured in CI"></a>
  <a href="https://github.com/kinglet-dev/core/releases/latest"><img src="https://raw.githubusercontent.com/kinglet-dev/core/badges/release.svg" alt="Latest release"></a>
  <a href="LICENSE"><img src="https://raw.githubusercontent.com/kinglet-dev/core/badges/licence.svg" alt="Licence: MIT"></a>
  <a href="go.mod"><img src="https://raw.githubusercontent.com/kinglet-dev/core/badges/go.svg" alt="Go version"></a>
</p>

# core

**The shared building blocks of Kinglet's maker tools:** a strict SVG reader that measures designs in millimetres, plane geometry, and the measurements behind [laserlint](https://github.com/kinglet-dev/laserlint)'s checks.

It's a Go library, standard library only. Nothing here touches the network.

[kinglet.dev](https://kinglet.dev) · [Changelog](CHANGELOG.md) · [Security](SECURITY.md)

> **Not yet stable.** core is at version 0, so its API can change between minor versions while the tools that use it settle. Pin a version in your `go.mod`.

## Example

How much line in [`examples/grid.svg`](examples/grid.svg), a grid of separate squares, would a laser burn twice?

```go
f, err := os.Open("grid.svg")
if err != nil {
	log.Fatal(err)
}
defer f.Close()

d, err := reader.Read(f) // millimetres, at the size saved in the file
if err != nil {
	log.Fatal(err)
}
r, err := crossing.Find(d.ScoredLines(0.005), crossing.Params{
	Merge: 0.1, Cell: 1, MaxSegments: 1_000_000, MaxEntries: 10_000_000,
})
if err != nil {
	log.Fatal(err)
}
fmt.Printf("%.1f × %.1f mm: %.1f mm of line burns twice, %d crossings\n",
	d.Width, d.Height, r.Stacked, len(r.Crossings))
```

```text
100.0 × 80.0 mm: 340.0 mm of line burns twice, 0 crossings
```

Every inner edge is drawn twice, once by each square beside it. A test runs this example, so the output above stays true.

## Packages

| Package | What it does |
|---|---|
| `adapters/svg/reader` | Reads an SVG into a design at its physical size, with limits on size, depth and element count. Refuses what it can't measure (`<text>`, `<image>`, `<use>`, CSS classes), with a fix for each. |
| `adapters/svg/syntax` | Parses SVG path data and transforms. |
| `adapters/svg/edit` | Replaces or removes chosen elements and keeps every other byte of the file as it was. |
| `adapters/files` | Opens a file or standard input within a size limit. |
| `adapters/buildinfo` | Works out the running version, for `--version`. |
| `domain/design` | A design: its size and shapes (with each shape's element and paint), and the lines a laser scores. |
| `domain/geom` | Points, affine transforms, paths, polylines; curves flattened within a tolerance. |
| `domain/units` | Lengths with units (`mm`, `cm`, `in`, `pt`, `px`). |
| `domain/crossing` | Where lines cross, and line lying on other line. |
| `domain/spacing` | Line running closer than a minimum distance to other line. |
| `domain/density` | The densest area of line. |
| `domain/enclose` | The smallest circle around a shape (Welzl's algorithm). |
| `domain/fill` | Whether a point is inside a filled shape (nonzero rule). |
| `domain/parallel` | Runs work on all processors. |

## Installation

```sh
go get github.com/kinglet-dev/core@latest
```

## Compatibility

Go 1.27 or later, on Windows, macOS and Linux (CI runs on all three). No dependencies outside the standard library; the tests use [Ginkgo](https://github.com/onsi/ginkgo) and [Gomega](https://github.com/onsi/gomega) (MIT).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Report security problems privately, as [SECURITY.md](SECURITY.md) describes.

## License

[MIT](LICENSE). The Kinglet name, logo and mascot are not covered by it; see [TRADEMARKS.md](TRADEMARKS.md).
