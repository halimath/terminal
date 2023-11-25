# termx

eXtended ANSI terminal support for Go(lang)

![CI Status][ci-img-url]
[![Package Doc][package-doc-img-url]][package-doc-url] 
[![Releases][release-img-url]][release-url]

[ci-img-url]: https://github.com/halimath/termx/workflows/CI/badge.svg
[package-doc-img-url]: https://img.shields.io/badge/GoDoc-Reference-blue.svg
[package-doc-url]: https://pkg.go.dev/github.com/halimath/termx
[release-img-url]: https://img.shields.io/github/v/release/halimath/termx.svg
[release-url]: https://github.com/halimath/termx/releases

`termx` is an extension of [`golang.org/x/term`](https:golang.org/x/term) that provides convenience functions
and types to ease creating applications that make intensive use of ANSI terminal features such as

* colored output
* raw input handling
* responsive rendering

# Installation

`termx` can be installed as a Go module

```shell
go get github.com/halimath/termx
```

# Usage

# Useful resources

A lot of information about the escape sequences and their treatment by different terminal applications has
been gathered from the following resources which i highly recommend to anyone interested in the internals
of `termx`:

* https://invisible-island.net/xterm/ctlseqs/ctlseqs.html
* https://chrisyeh96.github.io/2020/03/28/terminal-colors.html

# License

Copyright 2023 Alexander Metzner.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
