# OpenAPI Client Examples

The API specification is defined in [`openapi.yaml`](./openapi.yaml).

## Generate A Go Client (example)

```bash
openapi-generator-cli generate \
  -i ./11-flagship/01-opslane/docs/openapi.yaml \
  -g go \
  -o ./tmp/opslane-go-client
```

## Generate A TypeScript Client (example)

```bash
openapi-generator-cli generate \
  -i ./11-flagship/01-opslane/docs/openapi.yaml \
  -g typescript-fetch \
  -o ./tmp/opslane-ts-client
```
