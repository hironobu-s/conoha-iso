[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)  [![Build Status](https://travis-ci.org/hironobu-s/conoha-iso.svg?branch=master)](https://travis-ci.org/hironobu-s/conoha-iso) [![codebeat badge](https://codebeat.co/badges/792c6579-ec06-4841-a6e2-d49df29c0640)](https://codebeat.co/projects/github-com-hironobu-s-conoha-iso)

# ConoHa ISO

[English is here](README-en.md)

[ConoHa](https://www.conoha.jp/)にISOイメージのダウンロードリクエストを送ったり、VPSへのISOイメージの挿入、排出などが行える簡易ツールです。ConoHaはAPI経由でしかISOイメージを扱えませんが、このツールを使うとコマンドライン/Webブラウザから簡単に扱うことができます。

## インストール

以下の手順で実行ファイルをダウンロードしてください。以下のコマンドはカレントディレクトリにダウンロードしますが、使用頻度が高い場合はパスの通った場所に置いてください。

**Mac OSX**

```bash
curl -sL https://github.com/hironobu-s/conoha-iso/releases/download/current/conoha-iso-osx.amd64.gz | zcat > conoha-iso && chmod +x ./conoha-iso
```

**Linux(amd64)**

```bash
curl -sL https://github.com/hironobu-s/conoha-iso/releases/download/current/conoha-iso-linux.amd64.gz | zcat > conoha-iso && chmod +x ./conoha-iso
```

**Windows(amd64)**

[ZIP file](https://github.com/hironobu-s/conoha-iso/releases/download/current/conoha-iso.amd64.zip)


## Dockerで使う

Dockerイメージを用意してあります。環境変数などを渡す必要があるので、添付のスクリプト[docker-conoha-iso.sh](https://github.com/hironobu-s/conoha-iso/blob/master/docker-conoha-iso.sh)を使うと簡単です。

[hironobu/conoha-iso](https://hub.docker.com/r/hironobu/conoha-iso/)

## 自分でビルドする

上述したように配布された実行ファイルを使用することもできますが、自分でビルドすることもできます。

```bash
make
```

binディレクトリ以下に実行ファイルと配布アーカイブファイルが生成されます。


## はじめに(認証情報とリージョンの指定)

conoha-isoを実行するには、APIへの認証情報とリージョンの指定が必須となります。

API認証情報は「APIユーザ名」「APIパスワード」「テナント名 or テナントID」です。これらの情報は[ConoHaのコントロールパネル](https://manage.conoha.jp/API/)にあります。

リージョンはISOイメージを登録するリージョンで、tyo1, tyo2, sin1, sjc1の4つです(順に東京、シンガポール、アメリカ)。

これらはコマンドライン引数で渡す方法と環境変数に登録する方法が選べます。

**コマンドライン引数で渡す**

-u -p -n -t -rオプションを使います。テナント名とテナントIDは、どちらか一方を指定するだけで良いです。たとえばlistコマンドを実行する場合、以下のようになります。また、リージョンは指定しなかった場合、tyo1が使用されます。

テナント名を指定
```bash
./conoha-iso list -u [APIユーザ名] -p [APIパスワード] -n [テナント名] -r [リージョン]
```

テナントIDを指定
```bash
./conoha-iso list -u [APIユーザ名] -p [APIパスワード] -t [テナントID] -r [リージョン]
```

**環境変数で渡す**

API認証情報は環境変数経由で渡すこともできます。変数名は OS_USERNAME, OS_PASSWORD, OS_TENANT_NAME, OS_AUTH_URL, OS_TENANT_ID, OS_REGIONです。以下はbashの場合です。[^1]

```bash
export OS_USERNAME=[APIユーザ名]
export OS_PASSWORD=[APIパスワード]
export OS_TENANT_NAME=[テナント名]
export OS_AUTH_URL=[Identity Endpoint]
export OS_REGION_NAME=[リージョン]
```

[^1]: 互換性維持のためCONOHA_で始まる環境変数も使えます。

## 使い方

機能ごとにサブコマンドになっています。

### download

ConoHaにISOイメージをダウンロードするよう要求します。以前のConoHaのISOイメージアップロードに近い機能です。-iオプションでISOイメージのURLを指定してください。

```bash
./conoha-iso download -i http://stable.release.core-os.net/amd64-usr/current/coreos_production_iso_image.iso
```

### list

登録されているISOイメージの一覧を取得できます。ダウンロード要求が完了するには少し時間がかかりますので、downloadコマンド実行後にlistコマンドでチェックしてください。

```bash
./conoha-iso list
```

出力例

```
# ./conoha-iso list
[Image1]
Name:  alpine-mini-3.2.0-x86_64.iso
Url:   http://wiki.alpinelinux.org/cgi-bin/dl.cgi/v3.2/releases/x86_64/alpine-mini-3.2.0-x86_64.iso
Path:  /mnt/isos/repos/tenant_iso_data/6150e7c42bab40c59db53d415629841f/alpine-mini-3.2.0-x86_64.iso
Ctime: Wed May 27 04:30:45 2015
Size:  92329984

[Image2]
Name:  coreos_production_iso_image.iso
Url:   http://beta.release.core-os.net/amd64-usr/current/coreos_production_iso_image.iso
Path:  /mnt/isos/repos/tenant_iso_data/6150e7c42bab40c59db53d415629841f/coreos_production_iso_image.iso
Ctime: Thu May 28 02:03:18 2015
Size:  178257920
```

### insert

VPSにISOイメージを挿入します。実行するとイメージを挿入するVPSを選択するメニューと、挿入するISOイメージを選択するメニューが順に表示されるので、番号で選択してください。

```
# ./conoha-iso insert
[1] ***-***-***-***
[2] ***-***-***-***
Please select VPS no. [1-2]: 1

[1] alpine-mini-3.2.0-x86_64.iso
[2] coreos_production_iso_image.iso
Please select ISO no. [1-2]: 2
INFO[0039] ISO file was inserted and changed boot device.
```

### eject

VPSからISOイメージを排出します。
VPSを選択するメニューが表示されるので、番号で選択してください。

```
# ./conoha-iso eject
[1] ***-***-***-***
[2] ***-***-***-***
Please select VPS no. [1-2]: 1
INFO[0001] ISO file was ejected.
```

### server

WebブラウザからISOイメージの操作が行える管理コンソールを起動します。実行するとURLが表示されるので、Webブラウザを開き、アクセスしてください。デフォルトのURLは http://127.0.0.1:6543/ です。

```
$ ./conoha-iso server
Running on http://127.0.0.1:6543/
```

**-l** オプションでWebサーバが利用するアドレスとポート番号を指定することができます。

```
$ ./conoha-iso server -l 0.0.0.0:10000
Running on http://0.0.0.0:10000/
```

![conoha-iso-webui.png](conoha-iso-webui.png)

### ヘルプ

各サブコマンドは-hオプションをつけることでヘルプが出ます。

```bash
./conoha-iso -h
```

## ライセンス

MIT
