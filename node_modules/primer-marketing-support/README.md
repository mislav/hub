# Primer CSS Marketing Support

[![npm version](http://img.shields.io/npm/v/primer-marketing-support.svg)](https://www.npmjs.org/package/primer-marketing-support)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

> Support files are Sass variables, mixins, and functions that we import into different bases for use across components, objects, and utilities. Sharing these common properties across GitHub sites helps us to keep our styles more consistent.
>
> Most of the time to include these you'll only need to add `@import "./primer-marketing-support";` to the top of your bundle. If you want only a specific partial you can import them separately.

This repository is a module of the full [primer-css][primer] repository.

## Documentation

<!-- %docs
title: Variables
status: In review
-->

Documentation & refactor coming very soon

<!-- %enddocs -->

## Install

This repository is distributed with [npm][npm]. After [installing npm][install-npm], you can install `support` with this command.

```
$ npm install --save support
```

## Usage

The source files included are written in [Sass][sass] (`scss`) You can simply point your sass `include-path` at your `node_modules` directory and import it like this.

```scss
@import "primer-marketing-support/index.scss";
```

You can also import specific portions of the module by importing those partials from the `/lib/` folder. _Make sure you import any requirements along with the modules._

## License

MIT &copy; [GitHub](https://github.com/)

[primer]: https://github.com/primer/primer
[docs]: http://primercss.io/
[npm]: https://www.npmjs.com/
[install-npm]: https://docs.npmjs.com/getting-started/installing-node
[sass]: http://sass-lang.com/
