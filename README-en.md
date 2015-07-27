# ConoHa ISO

This is a simple tool that send the download request to APIs in [ConoHa](https://www.conoha.jp/). This also execute to insert or eject ISO image from VPS. You can handle only ISO image via APIs in ConoHa. but you can handle it from CLI more easily.

## Install

Please download the executable files by the following.

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


## Introduction

You need the authentication information such as API-Username, API-Password, Tenant-ID and Region to run conoha-iso. Region should be "tyo1", "sin1", or "sjc1".

These are on the ConoHa control-panel. 

How to pass these, You can select the way via command-line arguments, or also environment variables.

**Via command-line arguments**

You can use -p, -t, -r options.

```bash
./conoha-iso list -u [API-Username] -p [API-Password] -t [Tenant-ID] -r [Region]
```

**Via environment variables**

Also you can use CONOHA_USERNAME, CONOHA_PASSWORD, CONOHA_TENANT_ID and CONOHA_REGION. For bash script.

```bash
export CONOHA_USERNAME=[API Username]
export CONOHA_PASSWORD=[API Password]
export CONOHA_TENANT_ID=[Tenant ID]
export CONOHA_REGION=[Regiona]
```

## How to use

Sub-commands are provided for each function.

### list

Get ISO image list. You may run it after download sub-command.

```bash
./conoha-iso download -i http://stable.release.core-os.net/amd64-usr/current/coreos_production_iso_image.iso
```

Output:

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

Insert an ISO image to your VPS. If you run it, The menu will be displayed to select ISO image.

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

Eject an ISO image from your VPS.


```
# ./conoha-iso eject
[1] ***-***-***-***
[2] ***-***-***-***
Please select VPS no. [1-2]: 1
INFO[0001] ISO file was ejected.
```

## Help

All sub-command are accepted -h option to display the descriptions.

```bash
./conoha-iso -h
```

## License

BSD License
