# Primer CSS

[![npm version](http://img.shields.io/npm/v/primer-css.svg)](https://www.npmjs.org/package/primer-css)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

Primer is the CSS framework that powers GitHub's front-end design. Primer CSS includes 23 packages that are grouped into 3 core meta-packages for easy install. Each package and meta-package is independently versioned and distributed via npm, so it's easy to include all or part of Primer within your own project.

The Primer CSS repo is managed as a monorepo that is composed of many npm packages.

## Install

This repository is distributed with [npm][npm]. After [installing npm][install-npm], you can install `Primer-CSS` with this command.

```
$ npm install --save primer-css
```

## Usage

The source files included are written in [Sass][sass] (`scss`) You can simply point your sass `include-path` at your `node_modules` directory and import it like this.

```scss
@import "primer-css/index.scss";
```

You can also import specific portions of the module by importing those partials from the `/lib/` folder. _Make sure you import any requirements along with the modules._

## Build

For a compiled **CSS** version of this module, a npm script is included that will output a CSS version to `build/build.css` The built CSS file is also included in the npm package.

```
$ npm run build
```

## Documentation

You can read more about primer in the [docs][docs].

## License

[MIT](./LICENSE) &copy; [GitHub](https://github.com/)

[primer]: https://github.com/primer/primer
[docs]: http://primercss.io/
[npm]: https://www.npmjs.com/
[install-npm]: https://docs.npmjs.com/getting-started/installing-node
[sass]: http://sass-lang.com/
