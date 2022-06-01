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

### SQL User & Permissions

```sql
CREATE ROLE exampledb_read_only_role;
GRANT USAGE ON SCHEMA public TO exampledb_read_only_role;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO exampledb_read_only_role;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO exampledb_read_only_role;
CREATE USER read_only_username WITH PASSWORD 'supersecretpassword';
GRANT exampledb_read_only_role TO read_only_username;
GRANT pg_read_all_stats TO exampledb_read_only_role;
```

### Check definition

```yaml
api_version: core/v2
type: CheckConfig
metadata:
  labels:
    sensu.io/workflow: ci_action 
  name: metric-sensu-pg-metric-log
spec:
  runtime_assets:
    - sensu-pg-metric-log
  command:  sensu-pg-metric-log --pgurl postgres://user:pass@pghost:5432/postgres --enable --glhost graylog.example.com --glport 12249
  subscriptions:
    - worker-drone
  publish: true
  interval: 60
  output_metric_format: graphite_plaintext
  output_metric_handlers:
    - influxdbh
  round_robin: true
  proxy_entity_name: round_robin
  handlers:
    - notify_all
---
type: Asset
api_version: core/v2
metadata:
  name: sensu-pg-metric-log
  labels:
    sensu.io/workflow: ci_action
  annotations:
    io.sensu.bonsai.url: https://bonsai.sensu.io/assets/DoctorOgg/sensu-pg-metric-log
    io.sensu.bonsai.api_url: https://bonsai.sensu.io/api/v1/assets/DoctorOgg/sensu-pg-metric-log
    io.sensu.bonsai.tier: Community
    io.sensu.bonsai.version: 0.0.1
    io.sensu.bonsai.namespace: DoctorOgg
    io.sensu.bonsai.name: sensu-pg-metric-log
    io.sensu.bonsai.tags: ''
spec:
  builds:
  - url: https://assets.bonsai.sensu.io/848fab41fdae2c72373f3ea92091f20dfa22fbee/sensu-pg-metric-log_0.0.1_windows_amd64.tar.gz
    sha512: 9c1f2e318f9fd54623179d93cfce6001ea89520f87261aae9347e72e457defca4a09b7875502a24f33d0bffe2942c7bd143f96beaf820e07f902a3dabf5c0f4b
    filters:
    - entity.system.os == 'windows'
    - entity.system.arch == 'amd64'
  - url: https://assets.bonsai.sensu.io/848fab41fdae2c72373f3ea92091f20dfa22fbee/sensu-pg-metric-log_0.0.1_darwin_amd64.tar.gz
    sha512: e186f89b9bd9fedaae692647bc47e032ad5696389cf3857125fd4730be31be86d5c4eecc2f133447fc2f3ded12ae53c89afe08300eeb635e4ab9e8b82744863f
    filters:
    - entity.system.os == 'darwin'
    - entity.system.arch == 'amd64'
  - url: https://assets.bonsai.sensu.io/848fab41fdae2c72373f3ea92091f20dfa22fbee/sensu-pg-metric-log_0.0.1_linux_armv6.tar.gz
    sha512: ba60a3052d85d37bd9096f71fb81a92bdca42a3df927afa3fb2d168263b392976cdbf7b6c0b6b4de980d87e3ffd21a68f4974e56d4e7b64e2d750f872479342a
    filters:
    - entity.system.os == 'linux'
    - entity.system.arch == 'arm'
    - entity.system.arm_version == 6
  - url: https://assets.bonsai.sensu.io/848fab41fdae2c72373f3ea92091f20dfa22fbee/sensu-pg-metric-log_0.0.1_linux_armv7.tar.gz
    sha512: 140eb4f141fb2d4d91bb339769b5c88993fc23d26da976b0c7accd4c15455f4bb4b3614d766ac839973943e53145a317c302a41c0308f5c82be4ac5fb0b13793
    filters:
    - entity.system.os == 'linux'
    - entity.system.arch == 'arm'
    - entity.system.arm_version == 7
  - url: https://assets.bonsai.sensu.io/848fab41fdae2c72373f3ea92091f20dfa22fbee/sensu-pg-metric-log_0.0.1_linux_arm64.tar.gz
    sha512: 74df140dc00386808e3a1e62b021f0466786ea4c2ff48c42ec6413bc43faee5a3b00f45513af4cb859a94b4c9d933790fc64e7e95351519ad7893a34d41038e1
    filters:
    - entity.system.os == 'linux'
    - entity.system.arch == 'arm64'
  - url: https://assets.bonsai.sensu.io/848fab41fdae2c72373f3ea92091f20dfa22fbee/sensu-pg-metric-log_0.0.1_linux_amd64.tar.gz
    sha512: 12272a32925e8d5caf49fa9aa3e6905f252aac3d53c0c36788e7eda5c0c13633f12ed122d52f00c73b5b8acd5d235170de20955aea7cba8df87a89b8b1f78aa1
    filters:
    - entity.system.os == 'linux'
    - entity.system.arch == 'amd64'
```
