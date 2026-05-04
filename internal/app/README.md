# internal/app

Application-level runtime setup, configuration, and infrastructure.

## Responsibilities

- **Configuration**: `Config` struct and loading logic (from YAML/JSON or flags).
- **Initialization**: Setting up the data directory, default deck creation, and logging.
- **Logging**: Global logging configuration for the application.

## Key Symbols

- `Config`: Main application configuration.
- `SetupDataDir`: Ensures the workspace directory exists and contains necessary subfolders.
- `DefaultConfig`: Returns a sane default configuration for the app.
