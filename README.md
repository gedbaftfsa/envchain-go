# envchain-go

A CLI for managing per-project environment variable sets with encryption at rest.

## Installation

```bash
go install github.com/yourusername/envchain-go@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envchain-go/releases).

## Usage

**Store variables for a project namespace:**

```bash
envchain-go set myproject AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
envchain-go set myproject AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG
```

**Run a command with injected environment variables:**

```bash
envchain-go run myproject -- aws s3 ls
```

**List stored namespaces:**

```bash
envchain-go list
```

**Remove a variable or namespace:**

```bash
envchain-go unset myproject AWS_ACCESS_KEY_ID
envchain-go remove myproject
```

Variables are encrypted at rest using AES-GCM with a key derived from your system keychain or a master passphrase.

## Requirements

- Go 1.21+
- Linux, macOS, or Windows

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](LICENSE)