# mad

[中文](readme_zh.md)

[🗣 Talks](https://t.me/txthinking_talks)
[💬 Join](https://join.txthinking.com)
[🩸 Youtube](https://www.youtube.com/txthinking) 
[❤️ Sponsor](https://github.com/sponsors/txthinking)

Generate root CA and derivative certificate for any domains and any IPs

❤️ A project by [txthinking.com](https://www.txthinking.com)

### Install via [Nami](https://github.com/txthinking/nami)

    $ nami install github.com/txthinking/mad

### Usage

```
NAME:
   Mad - Generate root CA and derivative certificate for any domains and any IPs

USAGE:
   mad [global options] command [command options] [arguments...]

VERSION:
   20210401

AUTHOR:
   Cloud <cloud@txthinking.com>

COMMANDS:
   ca       Generate CA
   cert     Generate certificate
   install  Install ROOT CA
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

### Example

Generate root CA

```
$ sudo mad ca --ca ./ca.pem --key ./ca_key.pem --install
```

Generate cert for `localhost`

```
$ mad cert --ca ./ca.pem --ca_key ./ca_key.pem --cert ./localhost_cert.pem --key ./localhost_cert_key.pem --domain localhost
```

## License

Licensed under The MIT License
