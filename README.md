# envdiff

Compare `.env` files across staging and production configs and report missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
envdiff --base .env.staging --compare .env.production
```

**Example output:**

```
Missing in production:
  - DEBUG_MODE
  - FEATURE_FLAG_NEW_UI

Mismatched keys:
  - DATABASE_URL  (values differ)
  - LOG_LEVEL     staging="debug" production="warn"

2 missing, 2 mismatched
```

### Flags

| Flag        | Description                              | Default           |
|-------------|------------------------------------------|-------------------|
| `--base`    | Path to the base env file                | `.env.staging`    |
| `--compare` | Path to the env file to compare against  | `.env.production` |
| `--keys-only` | Only report missing keys, skip value diffs | `false`        |
| `--quiet`   | Exit with non-zero code on diff, no output | `false`          |

---

## Use in CI

```bash
envdiff --base .env.staging --compare .env.production --quiet || exit 1
```

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)