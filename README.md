# Peggo

<!-- markdownlint-disable MD041 -->

[![Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public.](https://img.shields.io/badge/repo%20status-WIP-yellow.svg?style=flat-square)](https://www.repostatus.org/#wip)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue?style=flat-square&logo=go)](https://godoc.org/github.com/umee-network/peggo)
[![Go Report Card](https://goreportcard.com/badge/github.com/umee-network/peggo?style=flat-square)](https://goreportcard.com/report/github.com/umee-network/peggo)
[![Version](https://img.shields.io/github/tag/umee-network/peggo.svg?style=flat-square)](https://github.com/umee-network/peggo/releases/latest)
[![License: Apache-2.0](https://img.shields.io/github/license/umee-network/peggo.svg?style=flat-square)](https://github.com/umee-network/peggo/blob/main/LICENSE)
[![Lines Of Code](https://img.shields.io/tokei/lines/github/umee-network/peggo?style=flat-square)](https://github.com/umee-network/peggo)
[![GitHub Super-Linter](https://img.shields.io/github/workflow/status/umee-network/peggo/Lint?style=flat-square&label=Lint)](https://github.com/marketplace/actions/super-linter)

Peggo is a Go implementation of the Peggy (Gravity Bridge) Orchestrator originally
implemented by [Injective Labs](https://github.com/InjectiveLabs/). Peggo itself
is a fork of the original Gravity Bridge Orchestrator implemented by [Althea](https://github.com/althea-net).


## Table of Contents

- [How it works](#how-it-works)
- [Dependencies](#dependencies)
- [Installation] (#installation)

## How it works

>  TLDR; It works by relaying messages from one chain to the other via relayers (validators).

## Dependencies

- [Go 1.17+](https://golang.org/dl/)

## Installation

To install the `peggo` binary:

```shell
$ make install
```