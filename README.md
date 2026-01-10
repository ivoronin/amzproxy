# amzproxy

HTTP reverse proxy that signs requests with AWS Signature Version 4

[![CI](https://github.com/ivoronin/amzproxy/actions/workflows/test.yml/badge.svg)](https://github.com/ivoronin/amzproxy/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/ivoronin/amzproxy)](https://github.com/ivoronin/amzproxy/releases)

[Overview](#overview) · [Features](#features) · [Installation](#installation) · [Usage](#usage) · [Configuration](#configuration) · [Requirements](#requirements) · [License](#license)

```bash
# Access OpenSearch Dashboards through IAM authentication
amzproxy -service es -region us-east-1 -host search-my-domain.us-east-1.es.amazonaws.com

# Then open in browser or use curl
curl http://localhost:8080/_dashboards/
```

## Overview

amzproxy acts as a reverse proxy that intercepts HTTP requests, signs them using AWS Signature Version 4, and forwards them to AWS services. It uses the AWS SDK credential provider chain to obtain credentials (environment variables, shared config files, EC2/ECS instance roles). All requests and responses are logged with timing information.

## Features

- Reverse proxy for any AWS service requiring SigV4 authentication
- Uses AWS SDK default credential provider chain
- Logs all requests with method, path, query string, source address, status code, and duration
- Listens on port 8080
- Forwards to HTTPS endpoints

## Installation

### GitHub Releases

Download from [Releases](https://github.com/ivoronin/amzproxy/releases).

### Homebrew

```bash
brew install ivoronin/tap/amzproxy
```

### Build from source

```bash
git clone https://github.com/ivoronin/amzproxy.git
cd amzproxy
make build
```

## Usage

### Command Line Flags

```bash
amzproxy -service <service> -region <region> -host <host>
```

| Flag | Description | Required |
|------|-------------|----------|
| `-service` | AWS service name (es, execute-api, s3, etc.) | Yes |
| `-region` | AWS region (us-east-1, eu-west-1, etc.) | Yes |
| `-host` | Target AWS service hostname | Yes |
| `-version` | Print version and exit | No |

### Examples

```bash
# Proxy to OpenSearch
amzproxy -service es -region us-east-1 -host search-my-domain.us-east-1.es.amazonaws.com

# Proxy to API Gateway
amzproxy -service execute-api -region eu-west-1 -host abc123.execute-api.eu-west-1.amazonaws.com

# Proxy to S3
amzproxy -service s3 -region us-west-2 -host my-bucket.s3.us-west-2.amazonaws.com
```

### How It Works

1. Start amzproxy with the target service, region, and host
2. Send HTTP requests to `http://localhost:8080/your/path`
3. amzproxy signs each request with SigV4 using your AWS credentials
4. Request is forwarded to `https://<host>/your/path`
5. Response is returned to the client

## Configuration

amzproxy uses the AWS SDK default credential provider chain. Credentials are resolved in this order:

1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`)
2. Shared credentials file (`~/.aws/credentials`)
3. Shared config file (`~/.aws/config`)
4. EC2 Instance Metadata Service (IMDS)
5. ECS container credentials

### Environment Variables

| Variable | Description |
|----------|-------------|
| `AWS_ACCESS_KEY_ID` | AWS access key |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key |
| `AWS_SESSION_TOKEN` | Session token (for temporary credentials) |
| `AWS_PROFILE` | Named profile to use from credentials file |
| `AWS_REGION` | Default region (not used by amzproxy, use `-region` flag) |

## Requirements

- AWS credentials with permissions to access the target service
- Network access to the target AWS service endpoint

## License

[GPL-3.0](LICENSE)
