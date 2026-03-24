# URL Shortener

Simple Go backend project for practicing clean project structure, config loading, and database integration.

## Run

```bash
go run ./cmd/url-shortener
```

## Config

The app reads config from `CONFIG_PATH`.
If `CONFIG_PATH` is not set, it uses `./config/local.yaml` by default.

Example:

```bash
CONFIG_PATH=./config/local.yaml go run ./cmd/url-shortener
```

## Project Structure

- `cmd/url-shortener` - application entrypoint
- `internal/config` - config loading
- `internal/storage` - storage layer (in progress)
- `config` - YAML configs for environments
