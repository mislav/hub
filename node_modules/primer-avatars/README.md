# Primer CSS Avatars

[![npm version](http://img.shields.io/npm/v/primer-avatars.svg)](https://www.npmjs.org/package/primer-avatars)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

> Avatars are images that users can set as their profile picture. On GitHub, they’re always going to be rounded squares. They can be custom photos, uploaded by users, or generated as Identicons as a placeholder.

This repository is a module of the full [primer-css][primer-css] repository.

## Install

This repository is distributed with [npm][npm]. After [installing npm][install-npm], you can install `primer-avatars` with this command.

```
$ npm install --save primer-avatars
```

## Usage

The source files included are written in [Sass][sass] (`scss`) You can simply point your sass `include-path` at your `node_modules` directory and import it like this.

```scss
@import "primer-avatars/index.scss";
```

You can also import specific portions of the module by importing those partials from the `/lib/` folder. _Make sure you import any requirements along with the modules._

## Build

For a compiled **css** version of this module, a npm script is included that will output a css version to `build/build.css` The built css file is also included in the npm package.

```
$ npm run build
```

## Documentation

<!-- %docs
title: Avatars
status: Stable
-->

Avatars are images that users can set as their profile picture. On GitHub, they're always going to be rounded squares. They can be custom photos, uploaded by users, or generated as Identicons as a placeholder.

{:toc}

## Basic example

Add `.avatar` to any `<img>` element to make it an avatar. This resets some key styles for alignment, address a Firefox image placeholder bug, and rounds the corners.

Be sure to set `width` and `height` attributes for maximum browser performance.

```html
<img class="avatar" alt="jonrohan" src="/jonrohan.png?v=3&s=144" width="72" height="72">
```

### Small avatars

We occasionally use smaller avatars. Anything less than `48px` wide should include the `.avatar-small` modifier class to reset the `border-radius` to a more appropriate level.

```html
<img class="avatar avatar-small" alt="jonrohan" src="/jonrohan.png?v=3&s=64" width="32" height="32">
```

### Parent-child avatars

When you need a larger parent avatar, and a smaller child one, overlaid slightly, use the parent-child classes.

```html
<div class="avatar-parent-child float-left">
  <img class="avatar" alt="jonrohan" src="/jonrohan.png?v=3&s=96" width="48" height="48">
  <img class="avatar avatar-child" alt="josh" src="/josh.png?v=3&s=40" width="20" height="20">
</div>
```

### Avatar stack

Stacked avatars can be used to show who is participating in thread when there is limited space available. When you hover over the stack, the avatars will reveal themselves. Optimally, you should put no more than 3 avatars in the stack.

```html
<span class="avatar-stack tooltipped tooltipped-s" aria-label="jonrohan, aaronshekey, and josh">
  <img alt="@jonrohan" class="avatar" height="39" alt="jonrohan" src="/jonrohan.png" width="39">
  <img alt="@aaronshekey" class="avatar" height="39" alt="aaronshekey" src="/aaronshekey.png" width="39">
  <img alt="@josh" class="avatar" height="39" alt="josh" src="/josh.png" width="39">
</span>
```

## Circle Badge

`.CircleBadge` allows for the display of badge-like icons or logos. They are used mostly with Octicons or partner integration icons.

`.CircleBadge` should have an `aria-label`, `title` (for a link), or an `alt` (for child `img` elements) attribute specified if there is no text-based alternative to describe it. If there is a text-based alternative or the icon has no semantic value, `aria-hidden="true"` or an empty `alt` attribute may be used.

### Small

```html
<a class="CircleBadge CircleBadge--small float-left mr-2" href="#small" title="Travis CI">
  <img src="<%= image_path "modules/site/travis-logo.png" %>"  class="CircleBadge-icon" alt="">
</a>
<a class="CircleBadge CircleBadge--small bg-yellow" title="Zap this!" href="#small">
  <%= octicon "zap",  :class => "CircleBadge-icon text-white" %>
</a>
```

### Medium

```html
<div class="CircleBadge CircleBadge--medium bg-gray-dark">
    <img src="<%= image_path "modules/site/travis-logo.png" %>"  alt="Travis CI" class="CircleBadge-icon">
</div>
```

### Large

```html
<div class="CircleBadge CircleBadge--large">
  <img src="<%= image_path "modules/site/travis-logo.png" %>"  alt="Travis CI" class="CircleBadge-icon">
</div>
```

### Dashed connection

For specific cases where two badges or more need to be shown as related or connected (such as integrations or specific product workflows), a `DashedConnection` class was created. Use utility classes to ensure badges are spaced correctly.

```html
<div class="DashedConnection">
  <ul class="d-flex list-style-none flex-justify-between" aria-label="A sample GitHub workflow">
    <li class="CircleBadge CircleBadge--small" aria-label="GitHub">
      <%= octicon "mark-github", :class => "width-full height-full" %>
    </li>

    <li class="CircleBadge CircleBadge--small" aria-label="Slack">
        <img src="<%= image_path "modules/site/logos/slack-logo.png" %>"  alt="" class="CircleBadge-icon">
    </li>

    <li class="CircleBadge CircleBadge--small" aria-label="Travis CI">
        <img src="<%= image_path "modules/site/travis-logo.png" %>"  alt="" class="CircleBadge-icon">
    </li>
  </ul>
</div>
```

<!-- %enddocs -->

## License

[MIT](./LICENSE) &copy; [GitHub](https://github.com/)

[primer-css]: https://github.com/primer/primer
[docs]: http://primercss.io/
[npm]: https://www.npmjs.com/
[install-npm]: https://docs.npmjs.com/getting-started/installing-node
[sass]: http://sass-lang.com/
