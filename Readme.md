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

This plugin adds support for [Vultr](https://www.vultr.com/) cloud instances to the `docker-machine` command line tool.

[![CircleCI](https://img.shields.io/circleci/project/janeczku/docker-machine-vultr.svg)](https://circleci.com/gh/janeczku/docker-machine-vultr/)

## Installation

Requirement: [Docker Machine 0.5.0-rc4](https://github.com/docker/machine)

Download the `docker-machine-driver-vultr` binary from the release page.
Extract the archive and copy the binary to a folder located in your `PATH` and make sure it's executable (e.g. `chmod +x /usr/local/bin/docker-machine-driver-vultr`).

## Usage instructions

Grab your API key from the [Vultr control panel](https://my.vultr.com/settings/) and pass that to `docker-machine create` with the `--vultr-api-key` option.

    $ docker-machine create --driver vultr --vultr-api-key=aa11bb22cc33 test-vps

Command line flags:

 - `--vultr-api-key`: **required** Your Vultr API key.
 - `--vultr-ssh-user`: SSH username for the new VPS
 - `--vultr-os-id`: Operating system ID to use (OSID). See [available OS IDs](https://www.vultr.com/api/#os_os_list).
 - `--vultr-region-id`: Region the VPS will be created in (DCID). See [available Region IDs](https://www.vultr.com/api/#regions_region_list).
 - `--vultr-plan-id`: Plan to use for this VPS (VPSPLANID). See [available Plan IDs](https://www.vultr.com/api/#plans_plan_list).
 - `--vultr-pxe-script`: PXE boot script ID. Requires custom OS (vultr-os-id=159)
 - `--vultr-ipv6`: Enable IPv6 support for the VPS. 
 - `--vultr-private-networking`: Enable private networking support for the VPS.
 - `--vultr-backups`: Enable automatic backups for the VPS.

      
### PXE deployment
OS installation on Vultr can take several minutes. If you need low ETAs for your machines i recommend to use iPXE boot method by choosing the `OSID 159`. By default this will boot the latest stable version of [RancherOS](http://rancher.com/rancher-os/) with persistent storage on the disk.

**Example for creating a new machine running RancherOS:**

    docker-machine create --driver vultr --vultr-api-key=aa11bb22cc33 --vultr-os-id=159 test-vps

 Alternatively you can use any PXE boot script that you created in your Vultr account panel by supplying it's ID via the `--vultr-pxe-script` flag.

 Environment variables and default values:

| CLI option                      | Environment variable         | Default                     |
|---------------------------------|------------------------------|-----------------------------|
| **`--vultr-api-key`**           | `VULTR_API_KEY`              | -                           |
| `--vultr-ssh-user`              | `VULTR_SSH_USER`             | `root`                      |
| `--vultr-os-id`                 | `VULTR_OS`                   | 160 (*Ubuntu 14.04 x64*)    |
| `--vultr-region-id`             | `VULTR_REGION`               | 1 (*New Jersey*)            |
| `--vultr-plan-id`               | `VULTR_PLAN`                 | 29 (*768 MB RAM,15 GB SSD*)|
| `--vultr-pxe-script`            | `VULTR_PXE_SCRIPT`           | -                           |
| `--vultr-ipv6`                  | `VULTR_IPV6`                 | `false`                     |
| `--vultr-private-networking`    | `VULTR_PRIVATE_NETWORKING`   | `false`                     |
| `--vultr-backups`               | `VULTR_BACKUPS`              | `false`                     |
     
