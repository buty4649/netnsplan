[English](README.md) | [Japanese](README_ja.md)

---

# netnsplan

`netnsplan` is a tool for easily setting up and managing complex network environments using Linux's network namespaces (netns). Network configurations are defined using YAML files, and the `netnsplan` command is used to apply and manage these configurations.

## Key Features

- **Configuration via YAML**: Network settings are defined in a YAML file, enabling high readability and manageability. The configuration file is similar to netplan, allowing for easy network configuration.
- **Network Namespace Management**: Supports the creation, configuration, and deletion of multiple network namespaces.
- **Flexible Network Configuration**: Supports configuration of physical devices, dummy interfaces (dummy devices), and Veth devices, as well as address assignment and routing settings.
- **Execution of Arbitrary Scripts**: Allows for the execution of any script on the network namespace after applying network settings.

## Usage

### Preparing the Configuration File

Network settings are defined in a YAML file. Here is an example:

```yaml
netns:
  ns1:
    ethernets:
      eth0:
        addresses:
          - 10.1.0.1/24
    dummy-devices:
      dummy0:
        addresses:
          - 10.2.0.1/24
    veth-devices:
      veth0:
        addresses:
          - 10.3.0.1/24
        peer:
          name: veth1
          netns: ns2
          addresses:
            - 10.3.0.2/24
```

In this configuration, two network namespaces, `ns1` and `ns2`, are created, each with different network interfaces configured.

### Executing Commands

To apply network namespaces and network settings based on the configuration file, execute the following command:

```bash
netnsplan apply -c config.yaml
```

### Deleting Network Namespaces

To delete the created network namespaces, execute the following command:

```bash
netnsplan destroy -c config.yaml
```

## Installation

### deb package

```sh
VERSION=$(wget -q -O- https://api.github.com/repos/buty4649/netnsplan/releases/latest | grep tag_name | cut -d '"' -f 4 | tr -d v)
case $(uname -m) in x86_64) ARCH=amd64;; aarch64) ARCH=arm64;; esac
wget https://github.com/buty4649/netnsplan/releases/download/v${VERSION}/netnsplan_${VERSION}_linux_${ARCH}.deb
sudo apt install ./netnsplan_${VERSION}_linux_${ARCH}.deb
```

### rpm package

```sh
VERSION=$(wget -q -O- https://api.github.com/repos/buty4649/netnsplan/releases/latest | grep tag_name | cut -d '"' -f 4 | tr -d v)
case $(uname -m) in x86_64) ARCH=amd64;; aarch64) ARCH=arm64;; esac
yum install https://github.com/buty4649/netnsplan/releases/download/v${VERSION}/netnsplan_${VERSION}_linux_${ARCH}.rpm
```

### Pre-built binary

```sh
VERSION=$(wget -q -O- https://api.github.com/repos/buty4649/netnsplan/releases/latest | grep tag_name | cut -d '"' -f 4 | tr -d v)
case $(uname -m) in x86_64) ARCH=amd64;; aarch64) ARCH=arm64;; esac
wget https://github.com/buty4649/netnsplan/releases/download/v${VERSION}/netnsplan_${VERSION}_linux_${ARCH}.tar.gz
tar xvf netnsplan_${VERSION}_linux_${ARCH}.tar.gz netnsplan
chmod +x netnsplan
sudo chown root: netnsplan
sudo mv netnsplan /usr/local/sbin/netnsplan
```

## License

See [LICENSE](LICENSE) © [buty4649](https://github.com/buty4649/)
