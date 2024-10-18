# Crawler SDK

A lightweight and efficient SDK designed to create web crawler applications while ensuring responsible behavior towards
target servers. The SDK achieves this by incorporating a proxy cache layer for all requests, minimizing direct
interactions with target servers and reducing the risk of overwhelming them.

## Features

- **Proxy Cache Layer**: Automatically caches HTTP requests using BadgerDB, reducing the number of requests sent to
  target servers.
- **Modular Design**: Split into different components for caching, data sources, and HTTP transport, allowing easy
  customization and extension.
- **Optimized for Performance**: Uses BadgerDB for fast, efficient storage and retrieval of cached responses.
- **Safe Crawling**: Designed with server safety in mind, making sure your application crawls responsibly without
  causing harm to the target servers.

## Getting Started

### Prerequisites

- Go 1.23+
- BadgerDB for caching

### Installation

Clone the repository and install the required dependencies:

```bash
git clone https://github.com/jictyvoo/radarada_crawler-sdk.git
cd crawler-sdk
go mod tidy
```

### Usage

Starting the Cache Proxy: Run the cache proxy to start intercepting and caching HTTP requests.

```bash
go run ./cmd/cacheproxy
```
