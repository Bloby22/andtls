# Contributing to andtls

Thank you for considering contributing to andtls!

## Getting Started

1. Fork the repository
2. Clone your fork
3. Create a feature branch: `git checkout -b feature/my-feature`
4. Make your changes
5. Run tests and linting:
   ```bash
   make fmt
   make vet
   ```
6. Commit your changes (see [Commit Guidelines](#commit-guidelines))
7. Push to your fork and open a Pull Request

## Development Setup

### Requirements

- **Go** 1.25.0+
- **ADB** (Android Debug Bridge) in PATH
- **Perl** 5.20+ (optional, for scripts)

### Build & Run

```bash
make build    # Build binary
make run      # Run directly via go run
make fmt      # Format code
make vet      # Run static analysis
```

### Project Structure

```
cmd/andtls/        # CLI entry point
internal/
  adb/             # ADB client and output parsing
  app/             # Application lifecycle
  config/          # Configuration
  device/          # Device model and USB watcher
  ui/              # Bubble Tea TUI
    components/    # UI components (table, details, help, etc.)
    styles/        # Colors and theme
scripts/           # Automation scripts (sh, pl)
```

## Commit Guidelines

- Use clear, descriptive commit messages
- Prefix commits with a category when appropriate:
  - `feat:` — new feature
  - `fix:` — bug fix
  - `docs:` — documentation changes
  - `refactor:` — code restructuring without behavior change
  - `chore:` — maintenance tasks

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Keep functions focused and small
- Add comments for non-obvious logic
- Use meaningful variable and function names

## Pull Request Process

1. Update documentation if your change affects usage
2. Ensure `make fmt` and `make vet` pass
3. Describe what your PR does and why
4. Reference any related issues

## Questions?

Open a discussion or issue on GitHub.
