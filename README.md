# mad

Generate root CA and derivative certificate for any domains and any IPs

❤️ A project by [txthinking.com](https://www.txthinking.com)

### Install via [Nami](https://github.com/txthinking/nami)

    nami install mad

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
mad ca --ca ./ca.pem --key ./ca_key.pem
```

Generate cert for `localhost`

```
mad cert --ca ./ca.pem --ca_key ./ca_key.pem --cert ./localhost_cert.pem --key ./localhost_cert_key.pem --domain localhost
```

## License

Licensed under The MIT License
