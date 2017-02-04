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

[![Latest Version](https://img.shields.io/github/release/janeczku/docker-machine-vultr.svg?maxAge=8600)][release]
[![Github All Releases](https://img.shields.io/github/downloads/janeczku/docker-machine-vultr/total.svg?maxAge=8600)]()
[![CircleCI](https://img.shields.io/circleci/project/janeczku/docker-machine-vultr/master.svg?maxAge=8600)][circleci]
[![License](https://img.shields.io/github/license/janeczku/docker-machine-vultr.svg?maxAge=8600)]()

[release]: https://github.com/janeczku/docker-machine-vultr/releases
[circleci]: https://circleci.com/gh/janeczku/docker-machine-vultr/

Provision [Vultr](https://www.vultr.com/) cloud instances using the `docker-machine` command line tool.

## Installation

Requirement: [Docker Machine >= 0.5.1](https://github.com/docker/machine)

Download the `docker-machine-driver-vultr` binary from the release page.
Extract the archive and copy the binary to a folder located in your `PATH` and make sure it's executable (e.g. `chmod +x /usr/local/bin/docker-machine-driver-vultr`).

## Usage instructions

Grab your API key from the [Vultr control panel](https://my.vultr.com/settings/) and pass that to `docker-machine create` with the `--vultr-api-key` option.

**Example for creating a new machine running RancherOS:**

    docker-machine create --driver vultr --vultr-api-key=abc123 rancheros-machine

**Example for creating a new machine running Ubuntu 14.04:**

    docker-machine create --driver vultr --vultr-api-key=abc123 --vultr-os-id=160 ubuntu-machine

Command line flags:

 - `--vultr-api-key`: **required** Your Vultr API key.
 - `--vultr-ssh-user`: SSH username.
 - `--vultr-region-id`: Region the VPS will be created in (DCID). See [available Region IDs](https://www.vultr.com/api/#regions_region_list).
 - `--vultr-plan-id`: Plan to use for this VPS (VPSPLANID). See [available Plan IDs](https://www.vultr.com/api/#plans_plan_list).
 - `--vultr-os-id`: Operating system ID to use (OSID). See [available OS IDs](https://www.vultr.com/api/#os_os_list).
 - `--vultr-ros-version`: RancherOS version to use if an OSID was not specified (e.g. 'v0.6.0', 'latest').
 - `--vultr-pxe-script`: PXE boot script ID. Requires the 'custom OS' ('--vultr-os-id=159')
 - `--vultr-ssh-key-id`: Use an existing SSH key in your Vultr account instead of generating a new one.
 - `--vultr-ipv6`: Enable IPv6 support for the VPS.
 - `--vultr-private-networking`: Enable private networking support for the VPS.
 - `--vultr-backups`: Enable automatic backups for the VPS.
 - `--vultr-userdata`: Path to file with cloud-init user-data.
 - `--vultr-snapshot-id`: Using snapshot to create host.
 - `--vultr-api-endpoint`: Override default Vultr API endpoint URL.

If the OS ID is not specified, the driver defaults to using [RancherOS](http://rancher.com/rancher-os/) as operating system for the instance.
A particular RancherOS version can be selected using the `--vultr-ros-version` flag.

### PXE deployment
You can boot a custom OS using a PXE boot script that you created in your Vultr account panel by supplying it's ID with the `--vultr-pxe-script` flag and setting `--vultr-os-id` to `159`.
The operating system must support Cloud-init and be configured to use the `ec2` datasource type.

 Environment variables and default values:

| CLI option                      | Environment variable         | Default                     |
|---------------------------------|------------------------------|-----------------------------|
| **`--vultr-api-key`**           | `VULTR_API_KEY`              | -                           |
| `--vultr-ssh-user`              | `VULTR_SSH_USER`             | `root`                      |
| `--vultr-region-id`             | `VULTR_REGION`               | 1 (*New Jersey*)            |
| `--vultr-plan-id`               | `VULTR_PLAN`                 | 29 (*768 MB RAM,15 GB SSD*) |
| `--vultr-os-id`                 | `VULTR_OS`                   | -                           |
| `--vultr-ros-version`           | `VULTR_ROS_VERSION`          | v0.5.0                      |
| `--vultr-pxe-script`            | `VULTR_PXE_SCRIPT`           | -                           |
| `--vultr-ssh-key-id`            | `VULTR_SSH_KEY`              | -                           |
| `--vultr-ipv6`                  | `VULTR_IPV6`                 | `false`                     |
| `--vultr-private-networking`    | `VULTR_PRIVATE_NETWORKING`   | `false`                     |
| `--vultr-backups`               | `VULTR_BACKUPS`              | `false`                     |
| `--vultr-userdata`              | `VULTR_USERDATA`             | -                           |
| `--vultr-snapshot-id`           | `VULTR_SNAPSHOT`             | -                           |
| `--vultr-api-endpoint`          | `VULTR_API_ENDPOINT`         | -                           |

### Find available plans for all Vultr locations

Check out [vultr-status.appspot.com](http://vultr-status.appspot.com) for a live listing of the available plans per region. Get the applicable `--vultr-region-id` and `--vultr-plan-id` parameters with the click of a button.

[![vultr-status website](vultr-status-screenshot.png?raw=true)](http://vultr-status.appspot.com)
