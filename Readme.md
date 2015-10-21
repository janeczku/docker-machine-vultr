<!--[metadata]>
+++
title = "Vultr"
description = "Vultr driver for docker machine"
keywords = ["machine, Vultr, driver, docker"]
[menu.main]
parent="smn_machine_drivers"
+++
<![end-metadata]-->

# Docker Machine driver plugin for Vultr

This plugin adds support to deploy [Vultr](https://www.vultr.com/) cloud instances to the `docker-machine` command line tool.

## Installation

Requirement: [Docker Machine](https://github.com/docker/machine)

Download the `docker-machine-vultr` binary from the release page.
Extract the archive and copy the binary to a folder located in your `PATH` and make sure it's executable (e.g. `chmod +x /usr/local/bin/docker-machine-vultr`).

## Usage instructions

Grab your API key from the [Vultr control panel](https://my.vultr.com/settings/) and pass that to `docker-machine create` with the `--vultr-api-key` option.

    $ docker-machine create --driver vultr --vultr-api-key=aa11bb22cc33 test-vps

Command line flags:

 - `--vultr-api-key`: **required** Your Vultr API key.
 - `--vultr-os-id`: Operating system ID to use (OSID). See [available OS IDs](https://www.vultr.com/api/#os_os_list).
 - `--vultr-region-id`: Region the VPS will be created in (DCID). See [available Region IDs](https://www.vultr.com/api/#regions_region_list).
 - `--vultr-plan-id`: Plan to use for this VPS (VPSPLANID). See [available Plan IDs](https://www.vultr.com/api/#plans_plan_list).
 - `--vultr-ipv6`: Enable IPv6 support for the VPS. 
 - `--vultr-private-networking`: Enable private networking support for the VPS.
 - `--vultr-backups`: Enable automatic backups for the VPS.

      
#### Recommendation
By default the driver will provision new VPS with an Ubuntu 14.04 x64 installation. Be aware that the deployment of Ubuntu instances on Vultr can take several minutes.
If you need faster ETAs for your machines i recommend to use [RancherOS](http://rancher.com/rancher-os/) as operating system by choosing `OSID` `159`. This will install the latest stable version of RancherOS via iPXE and should not take more than 30 seconds:

    --vultr-os-id=159

 Environment variables and default values:

| CLI option                      | Environment variable         | Default                |
|---------------------------------|------------------------------|------------------------|
| **`--vultr-api-key`**           | `VULTR_API_KEY`              | -                      |
| `--vultr-os-id`                 | `VULTR_OS`                   | 160 *Ubuntu 14.04 x64* |
| `--vultr-region-id`             | `VULTR_REGION`               | 1 *New Jersey*         |
| `--vultr-plan-id`               | `VULTR_PLAN`                 | 29 *768 MB RAM*        |
| `--vultr-ipv6`                  | `VULTR_IPV6`                 | `false`                |
| `--vultr-private-networking`    | `VULTR_PRIVATE_NETWORKING`   | `false`                |
| `--vultr-backups`               | `VULTR_BACKUPS`              | `false`                |
     
