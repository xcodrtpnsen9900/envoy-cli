# envoy-cli

A lightweight CLI for managing and switching between `.env` file profiles across projects.

---

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envoy-cli.git
cd envoy-cli && go build -o envoy .
```

---

## Usage

```bash
# Initialize envoy in your project
envoy init

# Save the current .env as a named profile
envoy save staging

# List available profiles
envoy list

# Switch to a profile
envoy use staging

# Show the active profile
envoy current
```

Profiles are stored in a `.envoy/` directory at the project root. You can safely commit this directory to version control — secrets stay in `.env`, which remains gitignored.

---

## Why envoy-cli?

Managing multiple environments (dev, staging, production) often means juggling `.env.staging`, `.env.production`, and manual copying. `envoy-cli` keeps your workflow clean with named, switchable profiles and a simple command interface.

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you'd like to change.

---

## License

[MIT](LICENSE)