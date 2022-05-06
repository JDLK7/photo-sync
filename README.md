- [Description](#description)
  - [Features](#features)
- [Usage](#usage)

# Description
Photo Sync is an agent that automatically synchronizes files when a new device is connected.
It was originally designed to copy media files from an SD card into an external drive. It was also designed to work in a Raspberry Pi
and includes a number of Ansible playbooks to bootstrap the server, configuring samba, and configuring the application itself.

## Features
- **Incremental backup:** Photo Sync will not copy already existing files. For instance, if you plug your camera SD card
and then use the same card to snap more pictures, it will only copy the new ones.
- **Notifications:** Photo Sync implements push notifications with [Pushover](https://pushover.net/). It will notify when
new device scans start/finish or when unexpected errors occur.

# Usage

1. Copy the configuration file
```shell
cp config.example.yaml
```
2. Edit the configuration file
3. Build the application
   1. On linux
    ```shell
    make build-linux
    ```
    2. On mac
    ```shell
    make build-darwin
    ```
4. Run the application
```shell
CONFIG_PATH=$pwd ./dist/linux/amd64/photo-sync
```
5. Plug a device
6. Mount the device
7. Start synching!
