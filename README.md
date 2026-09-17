# Go Package Repository

Go Package Repository is a small server for hosting [vanity Go import paths](https://go.dev/ref/mod#vcs-find). It serves the `go-import` and optional `go-source` HTML metadata consumed by the Go toolchain, while package records are managed through an authenticated JSON API and stored in SQLite.

For example, after registering `go.example.com/project`, users can run:

```console
go get go.example.com/project
```

The service returns metadata directing Go to the package's actual Git, Mercurial, Subversion, or Fossil repository.

The browser administration UI is documented in [ui/README.md](ui/README.md).

## Configuration

Configuration is read from environment variables. `APP_ADMIN_API_KEY` is the only required variable.

| Variable            | Required | Default          | Description                                                                                                                                                         |
|---------------------|----------|------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `APP_ADMIN_API_KEY` | Yes      | —                | Secret used to authenticate requests to `/admin/*`.                                                                                                                 |
| `APP_PORT`          | No       | `8080`           | HTTP listen port.                                                                                                                                                   |
| `APP_HOST_OVERRIDE` | No       | Request host     | Public host used when resolving an import path, such as `go.example.com`. This is useful behind a reverse proxy if it does not preserve the original `Host` header. |
| `APP_DATABASE_PATH` | No       | `./data.sqlite3` | Path to the SQLite database. The database and schema are created automatically.                                                                                     |

## Run with Docker

The image is published to the GitHub Container Registry:

```text
ghcr.io/rikudousage/go-pkg-repository
```

Release images have semantic-version tags such as `1.2.3`, `1.2`, and `1`. The `dev` tag follows the `master` branch. Pinning a release tag is recommended for deployments.

The container listens on port `8080` and runs as UID/GID `10001`. Create a writable directory for its SQLite database, then start it:

```console
mkdir -p data
sudo chown 10001:10001 data

docker run --name go-pkg-repository \
  --publish 8080:8080 \
  --env APP_ADMIN_API_KEY='replace-with-a-long-random-secret' \
  --env APP_HOST_OVERRIDE='go.example.com' \
  --env APP_DATABASE_PATH='/data/data.sqlite3' \
  --mount type=bind,source="$(pwd)/data",target=/data \
  ghcr.io/rikudousage/go-pkg-repository:dev
```

When running behind a reverse proxy, route the vanity domain to this service and preserve the `Host` header. Set `APP_HOST_OVERRIDE` when preserving it is not possible.

## Run from source

The project requires Go 1.27 and a C toolchain because its SQLite driver uses CGO.

```console
APP_ADMIN_API_KEY='replace-with-a-long-random-secret' \
APP_DATABASE_PATH='./data.sqlite3' \
go run .
```

Alternatively, with Nix installed:

```console
nix build .#app
APP_ADMIN_API_KEY='replace-with-a-long-random-secret' ./result/bin/go-pkg-repository
```

## Register a package

Administrative endpoints accept the API key either directly in the `Authorization` header or as a bearer token. Registering an existing import path updates its metadata.

```console
curl --fail-with-body \
  --request POST \
  --header 'Authorization: Bearer replace-with-a-long-random-secret' \
  --header 'Content-Type: application/json' \
  --data '{
    "import_path": "go.example.com/project",
    "vcs": "git",
    "repository_url": "https://github.com/example/project",
    "source_url": "https://github.com/example/project",
    "source_dir_url": "https://github.com/example/project/tree/main{/dir}",
    "source_file_url": "https://github.com/example/project/blob/main{/dir}/{file}#L{line}"
  }' \
  http://localhost:8080/admin/packages
```

The required fields are:

- `import_path`: the complete vanity import path, including its host;
- `vcs`: one of `git`, `hg`, `svn`, or `fossil`;
- `repository_url`: the repository clone URL.

The `source_url`, `source_dir_url`, and `source_file_url` fields are optional and populate the `go-source` metadata. Source URL templates may use the standard `{dir}`, `{file}`, and `{line}` placeholders.

The created record, including its numeric `id`, is returned as JSON with HTTP status `201 Created`.

## API

All administrative routes require the configured API key.

| Method | Path                            | Description                                                                     |
|--------|---------------------------------|---------------------------------------------------------------------------------|
| `POST` | `/admin/packages`               | Create a package or update one with the same import path.                       |
| `GET`  | `/admin/packages`               | List all packages.                                                              |
| `GET`  | `/admin/packages/{import-path}` | Get one package by its full import path.                                        |
| `GET`  | `/{import-path}`                | Serve Go discovery metadata; this public route does not require authentication. |

Examples:

```console
# List packages
curl --header 'Authorization: Bearer replace-with-a-long-random-secret' \
  http://localhost:8080/admin/packages

# Fetch one package
curl --header 'Authorization: Bearer replace-with-a-long-random-secret' \
  http://localhost:8080/admin/packages/go.example.com/project

# Inspect the public discovery metadata
curl 'http://localhost:8080/project?go-get=1'
```

This local example assumes `APP_HOST_OVERRIDE=go.example.com`. In a real deployment the public request normally arrives on the vanity host itself, so its path is likewise only the part after the hostname:

```console
curl 'https://go.example.com/project?go-get=1'
```

SQLite migrations run automatically whenever the service starts.
