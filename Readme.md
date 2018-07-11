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

You need to have installed [Docker Machine](https://github.com/docker/machine/releases) v0.5.1 or later.

Download the Vultr driver matching your platform from the [release page](https://github.com/janeczku/docker-machine-vultr/releases).

If you are running on **macOS**:

```console
$ curl -L https://github.com/janeczku/docker-machine-vultr/releases/download/v1.3.0/docker-machine-driver-vultr-`uname -s`-`uname -m` \
  -o /usr/local/bin/docker-machine-driver-vultr && chmod +x /usr/local/bin/docker-machine-driver-vultr
```

If you are running on **Linux**:

```console
$ curl -L https://github.com/janeczku/docker-machine-vultr/releases/download/v1.3.0/docker-machine-driver-vultr-`uname -s`-`uname -m` \
  -o /tmp/docker-machine-driver-vultr && chmod +x /tmp/docker-machine-driver-vultr &&
  sudo cp /tmp/docker-machine-driver-vultr /usr/local/bin/docker-machine-driver-vultr
```

Otherwise, download the binary from the release page directly.

## Usage instructions

Grab your API key from the [Vultr control panel](https://my.vultr.com/settings/) and pass that to `docker-machine create` with the `--vultr-api-key` option.

**Example for creating a new machine running RancherOS:**

    docker-machine create --driver vultr --vultr-api-key=abc123 rancheros-machine

**Example for creating a new machine running Ubuntu 16.04:**

    docker-machine create --driver vultr --vultr-api-key=abc123 --vultr-os-id=215 ubuntu-machine

Command line flags:

 - `--vultr-api-key`: **required** Your Vultr API key.
 - `--vultr-ssh-user`: SSH username.
 - `--vultr-region-id`: Region the VPS will be created in (DCID). See [available Region IDs](https://www.vultr.com/api/#regions_region_list).
 - `--vultr-plan-id`: Plan to use for this VPS (VPSPLANID). See [available Plan IDs](https://www.vultr.com/api/#plans_plan_list).
 - `--vultr-os-id`: Operating system ID to use (OSID). See [available OS IDs](https://www.vultr.com/api/#os_os_list).
 - `--vultr-app-id`: Vultr application ID to use with 'Application OS' ('--vultr-os-id=186').
 - `--vultr-ros-version`: RancherOS version to use if an OSID was not specified (e.g. 'v1.0.1', 'latest').
 - `--vultr-pxe-script`: PXE script ID. Requires the 'Custom OS' ('--vultr-os-id=159')
 - `--vultr-boot-script`: Boot script ID. Mutually exclusive of '--vultr-pxe-script'.
 - `--vultr-ssh-key-id`: Use an existing SSH key in your Vultr account instead of generating a new one.
 - `--vultr-ipv6`: Enable IPv6 support for the VPS.
 - `--vultr-private-networking`: Enable private networking support for the VPS.
 - `--vultr-backups`: Enable automatic backups for the VPS.
 - `--vultr-userdata`: Path to file with cloud-init user-data.
 - `--vultr-snapshot-id`: ID of an existing Snapshot in your Vultr account.
 - `--vultr-reserved-ip`: ID of a reserved IP in your Vultr account.
 - `--vultr-tag`: Tag to assign to the VPS.
 - `--vultr-firewall-group`: ID of existing firewall group to assign.
 - `--vultr-api-endpoint`: Override default Vultr API endpoint URL.

If the OS ID is not specified, [RancherOS](http://rancher.com/rancher-os/) will be used as operating system for the instance.
You can select a specific RancherOS version by specifying the `--vultr-ros-version` flag.

### PXE deployment
You can boot a custom OS using a PXE boot script that you created in your Vultr account panel by passing it's ID with the `--vultr-pxe-script` flag and setting `--vultr-os-id` to `159`.
The operating system must support cloud-init and be configured to use the `ec2` datasource type.

 Environment variables and default values:

| CLI option                      | Environment variable         | Default                     |
|---------------------------------|------------------------------|-----------------------------|
| **`--vultr-api-key`**           | `VULTR_API_KEY`              | -                           |
| `--vultr-ssh-user`              | `VULTR_SSH_USER`             | `root`                      |
| `--vultr-region-id`             | `VULTR_REGION`               | 1 (*New Jersey*)            |
| `--vultr-plan-id`               | `VULTR_PLAN`                 | 201 (*1024 MB, 25 GB SSD*)  |
| `--vultr-os-id`                 | `VULTR_OS`                   | -                           |
| `--vultr-ros-version`           | `VULTR_ROS_VERSION`          | v1.0.2                      |
| `--vultr-pxe-script`            | `VULTR_PXE_SCRIPT`           | -                           |
| `--vultr-boot-script`           | `VULTR_BOOT_SCRIPT`          | -                           |
| `--vultr-ssh-key-id`            | `VULTR_SSH_KEY`              | -                           |
| `--vultr-ipv6`                  | `VULTR_IPV6`                 | `false`                     |
| `--vultr-private-networking`    | `VULTR_PRIVATE_NETWORKING`   | `false`                     |
| `--vultr-backups`               | `VULTR_BACKUPS`              | `false`                     |
| `--vultr-userdata`              | `VULTR_USERDATA`             | -                           |
| `--vultr-snapshot-id`           | `VULTR_SNAPSHOT`             | -                           |
| `--vultr-reserved-ip`           | `VULTR_RESERVED_IP`          | -                           |
| `--vultr-tag`                   | `VULTR_TAG`                  | -                           |
| `--vultr-firewall-group`        | `VULTR_FIREWALL_GROUP`       | -                           |
| `--vultr-api-endpoint`          | `VULTR_API_ENDPOINT`         | -                           |

### Find available plans for all Vultr locations

Check out [vultr-status.appspot.com](http://vultr-status.appspot.com) for a live listing of the available plans per region. Get the corresponding `--vultr-region-id` and `--vultr-plan-id` parameters with the click of a button.

[![vultr-status website](vultr-status-screenshot.png?raw=true)](http://vultr-status.appspot.com)
