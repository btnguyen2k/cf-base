[![Actions Status](https://github.com/btnguyen2k/cf-base/workflows/ci/badge.svg)](https://github.com/btnguyen2k/cf-base/actions)
[![codecov](https://codecov.io/gh/btnguyen2k/cf-base/graph/badge.svg?token=HdvpHgjjvy)](https://codecov.io/gh/btnguyen2k/cf-base)
[![Release](https://img.shields.io/github/release/btnguyen2k/cf-base.svg?style=flat-square)](RELEASE-NOTES.md)

ContentFlow is a Content Management System for developers and technical writers.
Instead of relying on a GUI to create, publish, and update website content, it
uses developer-friendly tools and workflows: Git and CI/CD. ContentFlow is fast,
searchable, and requires no database backend.

This repository contains the models and schemas shared by other ContentFlow
components.

- [ContentFlow Engine](https://github.com/btnguyen2k/cf-engine)
- [ContentFlow CLI](https://github.com/btnguyen2k/cf-cli)

## Highlighted Features

- **Markdown** is a simple yet powerful markup language for creating formatted text. ContentFlow supports [GitHub Flavored Markdown](https://github.github.com/gfm/) as well as extensions such as Mathematical and Chemical formulas.
- **No Database** - website content is rendered from static data - which is Markdown text. Hence ContentFlow needs no database to run, and it is blazing fast.
- **Developer Friendly** - authoring website content is as similar as pushing code, making pull requests, builds, packages and deploying.
- **Fulltext Search** - website content is fulltext indexed and searchable.
- **Multi-languages** - multi-language website content is support, switching languages is on-the-fly.
