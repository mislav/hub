# Primer Marketing CSS Page Headers

[![npm version](http://img.shields.io/npm/v/primer-page-headers.svg)](https://www.npmjs.org/package/primer-page-headers)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

> Jumbotrons, heroes, and featured content sections for marketing websites at GitHub.

This repository is a module of the full [primer-css][primer] repository.

## Documentation

<!-- %docs
title: Page headers
status: In review
-->

Documentation & refactor coming very soon

<!-- %enddocs -->


## Install

This repository is distributed with [npm][npm]. After [installing npm][install-npm], you can install `primer-page-headers` with this command.

```
$ npm install --save primer-page-headers
```

## Usage

The source files included are written in [Sass][sass] (`scss`) You can simply point your sass `include-path` at your `node_modules` directory and import it like this.

```scss
@import "primer-page-headers/index.scss";
```

You can also import specific portions of the module by importing those partials from the `/lib/` folder. _Make sure you import any requirements along with the modules._

## Build

For a compiled **css** version of this module, a npm script is included that will output a css version to `build/build.css`

```
$ npm run build
```

## License

MIT &copy; [GitHub](https://github.com/)

[primer]: https://github.com/primer/primer
[primer-support]: https://github.com/primer/primer-support
[support]: https://github.com/primer/primer-support
[docs]: http://primercss.io/
[npm]: https://www.npmjs.com/
[install-npm]: https://docs.npmjs.com/getting-started/installing-node
[sass]: http://sass-lang.com/
