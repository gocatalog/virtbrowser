# virtbrowser
VirtBrowser

## Contents
- [Usage and Installation](#usage-and-installation)
- [Getting started](#getting-started)
- [Devlopment SetUp](#dev-setup)

## Usage and Installation
:warning: :warning: :warning:
This script will help you manage your Go application as a service on a Debian-based system, ensuring that Go 1.22.5 is installed.

To run the script directly from github you can use
`bash`
```

curl -sSL https://raw.githubusercontent.com/gocatalog/virtbrowser/HEAD/install.sh | bash -s setup

```
`wget`
```
wget -qO- https://raw.githubusercontent.com/gocatalog/virtbrowser/HEAD/install.sh | bash -s setup
```

:rocket: You can use the following commands to manage the service:

```
export VIRTBROWSER_REPO_DIR="/opt/virtbrowser"
cd VIRTBROWSER_REPO_DIR

# Update code
./install.sh -s pull

# Reabse
./install.sh -s rebase

# fresh Build
./install.sh -s build

#  Setup as virtbrowser service
./install.sh -s setup

# start virtbrowser service
./install.sh -s start

# start virtbrowser service
./install.sh -s stop

# restart virtbrowser service
./install.sh -s restart

```


## Getting started

1. TODO: how to use this project

## Dev Setup

1. TODO: how to dev setup