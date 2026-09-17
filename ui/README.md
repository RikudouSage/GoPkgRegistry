# Go Package Repository UI

The UI is served on port `4000`. Its container image is published as:

```text
ghcr.io/rikudousage/go-pkg-repository-ui
```

Configure it with environment variables at startup:

| Variable                 | Default | Description                                                                                                                                                                                                                    |
|--------------------------|---------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `API_URL`                | `/api`  | API URL used by browsers.                                                                                                                                                                                                      |
| `SSR_API_URL`            | —       | Internal, absolute API origin used while server-side rendering, for example `http://repository:8080`.                                                                                                                          |
| `SAME_ORIGIN`            | `true`  | Set to `false` when the API is on another origin; this enables cross-site cookie settings.                                                                                                                                     |
| `ALLOWED_HOSTS`          | —       | Comma-separated hostnames accepted by Angular SSR, for example `ui.example.com,ui.internal.example.com`. Do not include ports.                                                                                                 |
| `BASE_HREF`              | `/`     | Public base path for the UI.                                                                                                                                                                                                   |
| `NG_TRUST_PROXY_HEADERS` | —       | Comma-separated forwarded headers accepted from a trusted proxy, such as `x-forwarded-for`, `x-forwarded-host`, `x-forwarded-port`, `x-forwarded-prefix`, and `x-forwarded-proto`. Enable only headers set by a trusted proxy. |

For a reverse proxy that routes `/api` to the repository server, run:

```console
docker run --rm --publish 4000:4000 \
  --init \
  --env API_URL=/api \
  --env SSR_API_URL=http://repository:8080 \
  --env ALLOWED_HOSTS=ui.example.com \
  ghcr.io/rikudousage/go-pkg-repository-ui:dev
```
