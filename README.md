[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/DoctorOgg/sensu-pg-metric-log)
![goreleaser](https://github.com/DoctorOgg/sensu-check-http-go/workflows/goreleaser/badge.svg)

# sensu-pg-metric-log

## Table of Contents

- [Overview](#overview)
- [Files](#files)
- [Usage examples](#usage-examples)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)

## Overview

A Hacked up tool, to find slow responding queries, provide a count of them (Graphite format), and log the offenders to graylog

## Files

- sensu-pg-metric-log

## Usage examples

```bash
$ sensu-pg-metric-log  --pgurl postgres://user:password@127.0.0.1:7000/postgres  --enable --glhost graylog.example.com   --glport 12249
query 22 1654116654

```

Help:

```bash
$ $ ./sensu-pg-metric-log -h
A Hacked up tool, to find slow responding queries, provide a count of them, and log to graylog

Usage:
  sensu-pg-metric-log [flags]
  sensu-pg-metric-log [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number of this plugin

Flags:
      --enable          Log results to graylog
      --glhost string   hostname of the graylog server
      --glport int      port of the graylog server
  -h, --help            help for sensu-pg-metric-log
      --pgurl string    URL to the postgres database

Use "sensu-pg-metric-log [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add DoctorOgg/sensu-pg-metric-log
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/DoctorOgg/sensu-pg-metric-log].

### Check definition

TBD
