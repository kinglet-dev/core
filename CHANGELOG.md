# Changelog

All notable changes to core. Versions follow [Semantic Versioning](https://semver.org/). Until 1.0, minor versions may change the API; patch versions never do.

## 0.4.0 (unreleased)

### Added

- `units.Length.Millimeters` and `enclose.Circle.Center`, in American spelling like the rest of Kinglet.

### Deprecated

- `units.Length.Millimetres` and `enclose.Circle.Centre`. They still work and give the same values; a later release will remove them.

### Changed

- American spelling throughout the comments, tests and docs. The license badge is now `badges/license.svg`.

## 0.3.0 (2026-10-11)

### Added

- **Each shape keeps its path as the element draws it** (`Shape.Local`), in the element's own units, with the transform to millimeters (`Shape.ToMM`), so a tool can write a changed path back.
- **`edit.Change.Redraw`:** keeps an element and redraws it as a `<path>` with new path data. Its other attributes stay exactly as written, and its contents stay. Its own geometry attributes go, and so does editor data that would redraw the old shape. Path data with characters path data never has is refused.
- **`syntax.FormatPath`:** writes path data with absolute commands, rounded to a millionth of a unit.
- **Geometry:** `Part` cuts a segment between two parameters exactly (curves stay curves), `Matrix.Inverse` undoes a transform, and `Subpath.Reverse` runs a subpath the other way.

## 0.2.0 (2026-10-11)

### Added

- **Each shape records where it came from:** the number of its element in document order (the root `svg` is 0, and skipped elements still count), and its fill and stroke as written, after inheritance.
- **`adapters/svg/edit`:** replaces or removes elements by that number and keeps every other byte of the file, so ids, styles, comments and layout survive. It refuses numbers the document doesn't have, overlapping changes and malformed XML.

## 0.1.0 (2026-10-11)

The first release: the packages laserlint 0.1.2 used for reading and measuring designs, moved here unchanged so other Kinglet tools can share them.

- **SVG:** a strict reader that measures a design at its physical size, and a path and transform parser.
- **Geometry:** points, transforms, paths and polylines, and lengths with units.
- **Measurements:** crossings and stacked line, spacing between lines, density, smallest enclosing circles, and nonzero fill.
- **Tool support:** opening input within a size limit, and the running version.
