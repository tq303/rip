# rip

![version](https://img.shields.io/github/v/release/tq303/rip) ![build](https://github.com/tq303/rip/actions/workflows/release.yml/badge.svg) ![language](https://img.shields.io/badge/built%20with-Go-00ADD8) ![license](https://img.shields.io/badge/license-none-lightgrey)

Cross-platform CLI for flashing images to drives.

---

## Install

**macOS (Apple Silicon)**

```bash
curl -L https://github.com/tq303/rip/releases/latest/download/rip-darwin-arm64 -o /usr/local/bin/rip && chmod +x /usr/local/bin/rip
```

**macOS (Intel)**

```bash
curl -L https://github.com/tq303/rip/releases/latest/download/rip-darwin-amd64 -o /usr/local/bin/rip && chmod +x /usr/local/bin/rip
```

**Linux**

```bash
curl -L https://github.com/tq303/rip/releases/latest/download/rip-linux-amd64 -o /usr/local/bin/rip && chmod +x /usr/local/bin/rip
```

**Go**

```bash
go install github.com/tq303/rip@latest
```

**Local development**

```bash
make install
```

---

## Usage

```bash
rip [image]
```

Prompts you to select a drive, confirms before writing, then flashes the image. Accepts a local file or URL.

```bash
rip image.iso
rip image.img --buffer 8
rip https://example.com/image.iso
```

> Raw copy only — works for standard `.iso` and `.img` disk images. Does not handle images that require bootloader installation (e.g. Windows ISOs).

### Flags

| Flag             | Default | Description             |
| ---------------- | ------- | ----------------------- |
| `-b`, `--buffer` | `4`     | Write buffer size in MB |
