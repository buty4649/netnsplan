[English](README.md) | [日本語(Japanese)](README_ja.md)

---

# netnsplan

`netnsplan`は、Linuxのネットワーク名前空間（netns）を利用して複雑なネットワーク環境を簡単にセットアップ・管理するためのツールです。YAMLファイルを用いてネットワーク設定を定義し、`netnsplan`コマンドを通じてこれらの設定を適用・管理します。

## 主な特徴

- **YAMLによる設定**: ネットワークの設定をYAMLファイルで定義し、可読性の高い設定管理を実現します。netplanに似たコンフィグファイルで、簡単にネットワーク設定ができます。
- **ネットワーク名前空間の管理**: 複数のネットワーク名前空間の作成、設定、削除をサポートします。
- **柔軟なネットワーク設定**: 物理デバイス、ダミーデバイス(dummy)、vethデバイスの設定や、アドレス割り当て、ルーティング設定など、多様なネットワーク設定に対応します。
- **任意のスクリプトの実行**: ネットワーク設定適用後に任意のスクリプトをネットワーク名前空間上で実行できます。

## 使用方法

### 設定ファイルの準備

ネットワーク設定はYAMLファイルによって定義します。以下はその例です：

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

この設定では、`ns1`と`ns2`という二つのネットワーク名前空間を作成し、それぞれに異なるネットワークインターフェイスを設定しています。

### コマンドの実行

設定ファイルを元にネットワーク名前空間とネットワーク設定を適用するには、以下のコマンドを実行します：

```bash
netnsplan apply -c config.yaml
```

### ネットワーク名前空間の削除

作成したネットワーク名前空間を削除するには、以下のコマンドを実行します：

```bash
netnsplan destroy -c config.yaml
```

## インストール

### debパッケージ

```sh
VERSION=$(wget -q -O- https://api.github.com/repos/buty4649/netnsplan/releases/latest | grep tag_name | cut -d '"' -f 4 | tr -d v)
case $(uname -i) in x86_64) ARCH=amd64;; aarch64) ARCH=arm64;; esac
wget https://github.com/buty4649/netnsplan/releases/download/v${VERSION}/netnsplan_${VERSION}_linux_${ARCH}.deb
sudo apt install ./netnsplan_${VERSION}_linux_${ARCH}.deb
```

### rpmパッケージ

```sh
VERSION=$(wget -q -O- https://api.github.com/repos/buty4649/netnsplan/releases/latest | grep tag_name | cut -d '"' -f 4 | tr -d v)
case $(uname -i) in x86_64) ARCH=amd64;; aarch64) ARCH=arm64;; esac
yum install https://github.com/buty4649/netnsplan/releases/download/v${VERSION}/netnsplan_${VERSION}_linux_${ARCH}.rpm
```

### ビルド済みバイナリ

```sh
VERSION=$(wget -q -O- https://api.github.com/repos/buty4649/netnsplan/releases/latest | grep tag_name | cut -d '"' -f 4 | tr -d v)
case $(uname -i) in x86_64) ARCH=amd64;; aarch64) ARCH=arm64;; esac
wget https://github.com/buty4649/netnsplan/releases/download/v${VERSION}/netnsplan_${VERSION}_linux_${ARCH}.tar.gz
tar xvf netnsplan_${VERSION}_linux_${ARCH}.tar.gz netnsplan
chmod +x netnsplan
sudo chown root: netnsplan
sudo mv netnsplan /usr/local/sbin/netnsplan
```

## ライセンス

See [LICENSE](LICENSE) © [buty4649](https://github.com/buty4649/)
