# Primer marketing

[![npm version](http://img.shields.io/npm/v/primer-marketing.svg)](https://www.npmjs.org/package/primer-marketing)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

> Primer marketing is one of 3 meta-packages that belong to the Primer CSS framework. Primer marketing contains packages that are used on GitHub marketing websites.

This repository is a compilation of [several CSS packages](https://github.com/primer/primer-css). You can break it down into smaller sections using npm.

## Install

This repository is distributed with [npm][npm]. After [installing npm][install-npm], you can install `primer-css` with this command. **You must have npm >=3.0 installed to be able to satisfy the dependency paths**

```
$ npm install --save primer-marketing
```

## Usage

The source files included are written in [Sass][sass] (`scss`) You can simply point your sass `include-path` at your `node_modules` directory and import it like this.

```scss
@import "primer-marketing/index.scss";
```

You can also import specific portions of the module by importing those partials from the `/lib/` folder. _Make sure you import any requirements along with the modules._

## Build

For a compiled **css** version of this module, a npm script is included that will output a CSS version to `build/build.css` The built CSS file is also included in the npm package.

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
