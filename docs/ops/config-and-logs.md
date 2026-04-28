# Config And Logs

## Data Directory

Default data directory:

```text
os.UserConfigDir()/deutsch-tui
```

Override with:

```sh
go run ./cmd/deutsch-tui --data-dir /path/to/data
```

## Config File

The app writes `config.json` in the data directory when missing.

Current fields:

- `theme`: default `system`
- `keymap`: default `default`
- `ai_provider`: default `disabled`
- `log_level`: default `info`

## Logs

The app writes `deutsch-tui.log` in the data directory. Logs are local only and should stay out of Git.
