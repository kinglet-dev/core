# Changelog

All notable changes to core. Versions follow [Semantic Versioning](https://semver.org/). Until 1.0, minor versions may change the API; patch versions never do.

## 0.1.0 (unreleased)

The first release: the packages laserlint 0.1.2 used for reading and measuring designs, moved here unchanged so other Kinglet tools can share them.

- **SVG:** a strict reader that measures a design at its physical size, and a path and transform parser.
- **Geometry:** points, transforms, paths and polylines, and lengths with units.
- **Measurements:** crossings and stacked line, spacing between lines, density, smallest enclosing circles, and nonzero fill.
- **Tool support:** opening input within a size limit, and the running version.
