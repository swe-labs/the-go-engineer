# Track CFG: Configuration

## Mission

Master application configuration and startup discipline. Learn how to manage environment variables, configuration files (YAML/JSON), and command-line flags. Understand the **12-Factor App** principles for configuration and learn how to validate your application state **on boot** to prevent runtime failures.

## Stage Ownership

This track belongs to [10 Production Operations](../README.md).

## Track Map

| ID | Type | Surface | Mission | Requires |
| --- | --- | --- | --- | --- |
| `CFG.1` | Lesson | [Environment Variables](./1-environment-variables) | Master the `os.Getenv` and "Dotenv" patterns. | entry |
| `CFG.2` | Lesson | [Config Files](./2-configuration-files) | Use YAML/JSON for complex hierarchical settings. | `CFG.1` |
| `CFG.3` | Lesson | [Flag Parsing](./3-flag-parsing) | Master the standard library `flag` package. | `CFG.2` |
| `CFG.4` | Lesson | [12-Factor Principles](./4-twelve-factor-principles) | Learn the industry standard for cloud-native config. | `CFG.3` |
| `CFG.5` | Exercise | [Validation on Boot](./5-config-validation-on-boot) | Fail fast if the configuration is invalid. | `CFG.4` |

## Why This Track Matters

In production, your application will run in multiple environments (Development, Staging, Production).

1. **Portability**: Proper configuration allows the *exact same binary* to run in different environments just by changing external settings.
2. **Security**: Secrets (API keys, DB passwords) must be injected via configuration, never hardcoded in the source.
3. **Reliability**: Validating configuration on startup ensures that if a database URL is missing, the application fails immediately instead of waiting for a user request to fail minutes later.

## Next Step

After mastering configuration, learn how to monitor your application's health and performance. Continue to [Track OPS: Observability](../05-observability).
