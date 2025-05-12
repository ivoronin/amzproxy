# AWS SigV4 Proxy
![GitHub release (with filter)](https://img.shields.io/github/v/release/ivoronin/amzproxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/ivoronin/amzproxy)](https://goreportcard.com/report/github.com/ivoronin/amzproxy)
![GitHub last commit (branch)](https://img.shields.io/github/last-commit/ivoronin/amzproxy/main)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/ivoronin/amzproxy/goreleaser.yml)
![GitHub top language](https://img.shields.io/github/languages/top/ivoronin/amzproxy)

A lightweight HTTP proxy that transparently signs requests using AWS Signature Version 4. This tool allows you to forward HTTP requests to AWS services while automatically applying proper authentication headers using your default AWS credentials.

## 🚀 Features

- Acts as a reverse proxy for AWS services
- Automatically signs all forwarded requests with SigV4
- Supports any AWS service with minimal configuration
- Uses default AWS credential provider chain (env, config files, EC2/ECS roles)
- Logs incoming requests for observability

## 🔧 Usage

### Command Line Arguments

```bash
-service string   # AWS service name (e.g., execute-api, es, s3)
-region string    # AWS region (e.g., us-east-1)
-host string      # AWS service host (e.g., search-my-domain.us-east-1.es.amazonaws.com)
```

### Example

```bash
amzproxy -service es -region us-east-1 -host search-my-domain.us-east-1.es.amazonaws.com
```

Then send requests to:

```
http://localhost:8080/_dashboards/
```

The proxy will:

1.	Sign the request using AWS SigV4.
2.	Forward it to https://search-my-domain.us-east-1.es.amazonaws.com/_dashboards/.
3.	Return the response to the client.
