# aikido-sdk-go
Go SDK for the Aikido public API.

## Setup

1. Setup a new API client on Aikido via https://app.aikido.dev/settings/integrations/api/aikido/rest and assign `issues:read` permissions.

2. Create a configuration file:

```yaml
log:
  level: info

aikido:
  client_id: foo
  client_secret: bar
```

3. Run it with `aikido -config=config.yml`

## Compilation

Just run `make`.