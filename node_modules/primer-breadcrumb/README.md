# Primer Marketing CSS Breadcrumb Navigation

[![npm version](http://img.shields.io/npm/v/primer-breadcrumb.svg)](https://www.npmjs.org/package/primer-breadcrumb)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

> Breadcrumb navigation for GitHub's marketing pages with parents / grandparents.

This repository is a module of the full [primer-css][primer] repository.

## Documentation

<!-- %docs
title: Breadcrumbs
status: Stable
-->

Breadcrumbs are used to show taxonomical context on pages that are many levels deep in a site’s hierarchy. Breadcrumbs show and link to parent, grandparent, and sometimes great-grandparent pages. Breadcrumbs are most appropriate on pages that:

- Are many levels deep on a site
- Do not have a section-level navigation
- May need the ability to quickly go back to the previous (parent) page

#### Usage

```html
<nav aria-label="Breadcrumb">
  <ol>
    <li class="breadcrumb-item text-small"><a href="/business">Business</a></li>
    <li class="breadcrumb-item text-small"><a href="/business/customers">Customers</a></li>
    <li class="breadcrumb-item breadcrumb-item-selected text-small text-gray" aria-current="page">MailChimp</li>
  </ol>
</nav>
```

<!-- %enddocs -->

## License

MIT &copy; [GitHub](https://github.com/)

[primer]: https://github.com/primer/primer
[primer-support]: https://github.com/primer/primer-support
[support]: https://github.com/primer/primer-support
[docs]: http://primercss.io/
[npm]: https://www.npmjs.com/
[install-npm]: https://docs.npmjs.com/getting-started/installing-node
[sass]: http://sass-lang.com/
