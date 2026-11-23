# Platform setup — Windows, macOS, Linux

One codebase. Three native targets. Nothing here is a code change.

## Requirements

| | Windows | macOS | Linux |
|---|---|---|---|
| Go | 1.26+ | 1.26+ | 1.26+ |
| Node | 20+ | 20+ | 20+ |
| Wails CLI | `go install github.com/wailsapp/wails/v3/cmd/wails3@latest` | same | same |
| WebView | WebView2 (preinstalled on Win11) | WKWebView (OS) | GTK4 + WebKitGTK 6.0 |
| C toolchain | not needed (`CGO_ENABLED=0`) | Xcode CLT | `build-essential` |

### Linux packages
```bash
# Debian / Ubuntu
sudo apt install libgtk-4-dev libwebkitgtk-6.0-dev libsoup-3.0-dev build-essential pkg-config
# Arch
sudo pacman -S gtk4 webkitgtk-6.0 libsoup3 base-devel
# Fedora
sudo dnf install gtk4-devel webkitgtk6.0-devel libsoup3-devel gcc
```

### macOS
```bash
xcode-select --install
```

### Windows
```powershell
winget install GoLang.Go OpenJS.NodeJS
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```
`wails3 doctor` reports anything missing.

## Everyday commands

```bash
make setup      # toolchain check + deps
make dev        # hot reload, native window (needs the build/ assets below)
make test       # go tests + typecheck (no webview needed)
make build      # binary for the host OS into bin/ (needs the build/ assets)
make package    # native installer for the host OS (needs the build/ assets)
```

### Build assets

`make dev`, `make build` and `make package` read `build/config.yml` and `build/*/Taskfile.yml`. Those
files come from the Wails v3 template and are not committed yet. Until they are, use the headless
commands (`make test`, `make build-server`, `make e2e`). Scaffold them with:

```bash
wails3 init -n flow -t react-ts -d /tmp/flow-template
cp -r /tmp/flow-template/build .
```

Keep this repo's `cmd/`, `internal/` and `frontend/` — copy only `build/`.

## Cross-compilation

Wails v3 cross-compiles through Docker, one image per target (~800 MB, one-time). The Makefile
targets exist, but they depend on the `build/` assets above:

```bash
make cross-linux    # from macOS or Windows
make cross-win
make cross-mac      # unsigned binary; sign on a Mac before distributing
```

Notes that will bite you:
- **macOS targets need macOS for signing and notarisation.** Unsigned `cross-mac` output runs
  locally but Gatekeeper will block it on another machine.
- **Windows needs no cgo** (`CGO_ENABLED=0`), so Windows builds work from any host.
- **Linux needs cgo** for GTK, so build Linux on Linux (or in the Docker image).
- Server mode (`make build-server`) is pure Go and cross-compiles anywhere.

## Verification matrix before a release

| Check | Command | Where |
|---|---|---|
| unit + domain tests | `go test -tags server ./...` | any OS |
| parity e2e | `make e2e` | any OS (headless browser) |
| window opens, canvas renders | `make dev` | each OS once |
| installer produced | `make package` | each OS |

## Data location (SQLite `app.db`)

| OS | Path |
|---|---|
| macOS | `~/Library/Application Support/flow/app.db` |
| Linux | `~/.config/flow/app.db` |
| Windows | `%AppData%\flow\app.db` |

Path comes from `internal/platform.DataDir`. Never hardcode it elsewhere.
