# Matrix Bumble Bridge

## Project Layout

- `main.go`: bridge entrypoint and startup wiring.
- `bumble/`: Bumble API client, connector, commands, queue/cache logic.
- `repo/`: database models and data access helpers.
- `config/`: config structs and embedded example config.
- `example-config.yaml`: copy this to local `config.yaml`.
- `Makefile`: helper targets (`run`, `proxy`, `docker-build`).

## Quickstart

1. Install Go and required native deps (`libolm` is required by mautrix crypto).
2. Create local config:

```bash
cp example-config.yaml config.yaml
```

3. Set required env vars:

```bash
export LOCALPART=bumble
export DATABASE_PATH=./matrix-bumble.sqlite
```

4. Generate appservice registration:

```bash
go run . -g -c config.yaml -r registration.yaml
```

5. Configure `config.yaml`:
- Set `homeserver.address` and `homeserver.domain`.
- Copy `appservice.id`, `as_token`, `hs_token` from generated `registration.yaml`.
- Set `bridge.permissions` for your Matrix user.

6. Install `registration.yaml` on your homeserver and restart the homeserver.
7. Run the bridge:

```bash
go run . -c config.yaml -r registration.yaml
```

## Login Flow (Matrix Side)

In a management room with the bridge bot:

1. Run `help` to see commands.
2. Run `login`.
3. Paste the full Bumble web `curl` command when prompted.

The bridge extracts cookies (`aid`, `HDR-X-User-id`, `session`, `device_id`, etc.) and starts syncing conversations.

## Notes

- Local runtime files are ignored by default (`config.yaml`, `registration.yaml`, `registration_*.yaml`, `*.sqlite`, `logs/`).
- `ANTHROPIC_API_KEY` is read from env when constructing the AI client. If unset, AI-dependent features will fail if executed.

## Development

Run tests:

```bash
go test ./...
```

Show CLI options:

```bash
go run . -h
```
