# Primer Marketing CSS Typography

[![npm version](http://img.shields.io/npm/v/primer-marketing-type.svg)](https://www.npmjs.org/package/primer-marketing-type)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

> Flash messages, or alerts, inform users of successful or pending actions. Use them sparingly. Don’t show more than one at a time.

This repository is a module of the full [primer-css][primer] repository.

## Documentation

<!-- %docs
title: Typography
status: New Release
-->

The typography for our marketing pages differs slightly from what is in Primer's core--it is responsive, on a slightly different scale, and headlines are in a different font (Roboto).


## Heading Utilities

Use `.alt-h1` – `.alt-h6` to change an element's font, size, and weight on marketing pages.

```html

<p class="alt-h1">Pizza 1</p>
<p class="alt-h2">Pizza 2</p>
<p class="alt-h3">Pizza 3</p>
<p class="alt-h4">Pizza 4</p>
<p class="alt-h5">Pizza 5</p>
<p class="alt-h6">Pizza 6</p>

```

## Typographic Utilities

These utilities are meant to be used in addition to Primer's core utilities.

```html

<p class="alt-lead text-gray">I'm a lead paragraph. Bacon ipsum dolor amet tri-tip chicken kielbasa, cow swine beef corned beef ground round prosciutto hamburger porchetta sausage alcatra tail.</p>

<p class="text-gray">I'm a regular paragraph, except I'm gray. Bacon ipsum dolor amet tri-tip chicken kielbasa, cow swine beef corned beef ground round prosciutto hamburger porchetta sausage alcatra tail.</p>

<p class="text-gray alt-text-small">And I'm a tiny lil' baby paragraph. Bacon ipsum dolor amet tri-tip chicken kielbasa, cow swine beef corned beef ground round prosciutto hamburger porchetta sausage alcatra tail.</p>

<p class="alt-mono-font text-gray">Occasionally we want to use a mono font, there is a utility class for that.</p>

<p class="pullquote">I'm a pullquote. Someone said these words in real life, and now they're on the internet</p>

```

<!-- %enddocs -->

## Install

This repository is distributed with [npm][npm]. After [installing npm][install-npm], you can install `primer-marketing-typography` with this command.

```
$ npm install --save primer-marketing-typography
```

## Usage

The source files included are written in [Sass][sass] (`scss`) You can simply point your sass `include-path` at your `node_modules` directory and import it like this.

```scss
@import "primer-marketing-typography/index.scss";
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
