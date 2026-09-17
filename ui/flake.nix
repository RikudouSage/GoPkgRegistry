{
  description = "UI development environment and container image";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];
      devSystems = supportedSystems ++ [ "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      forAllDevSystems = nixpkgs.lib.genAttrs devSystems;
      versionFromEnv = builtins.getEnv "PACKAGE_VERSION";
      packageVersion = if versionFromEnv == "" then "0.0.0" else versionFromEnv;
      dockerTag = if versionFromEnv == "" then "dev" else versionFromEnv;
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          yarn = pkgs.yarn.override { nodejs = pkgs.nodejs_26; };

          entrypoint = pkgs.writeShellScriptBin "go-pkg-repository-ui-entrypoint" ''
            set -eu

            runtime_variables_dir=/app/browser
            runtime_variables_template="$runtime_variables_dir/runtime-variables.template.js"
            runtime_variables_file="$runtime_variables_dir/runtime-variables.js"
            cp "$runtime_variables_template" "$runtime_variables_file"
            chmod u+w "$runtime_variables_file"

            if [ "''${ALLOWED_HOSTS+x}" = x ]; then
              export NG_ALLOWED_HOSTS="$ALLOWED_HOSTS"
            fi

            ${pkgs.nodejs_26}/bin/node -e '
              const fs = require("node:fs");
              const file = process.argv[1];

              for (const name of ["API_URL", "SAME_ORIGIN", "SSR_API_URL"]) {
                if (Object.hasOwn(process.env, name)) {
                  fs.appendFileSync(file, `\ndefine(''${JSON.stringify(name)}, ''${JSON.stringify(process.env[name])});\n`);
                }
              }
            ' "$runtime_variables_file"

            exec ${pkgs.nodejs_26}/bin/node /app/server/server.mjs "$@"
          '';

          app = pkgs.stdenv.mkDerivation {
            pname = "go-pkg-repository-ui";
            version = packageVersion;
            src = self;

            nativeBuildInputs = [
              pkgs.nodejs_26
              yarn
              pkgs.yarnConfigHook
            ];

            yarnOfflineCache = pkgs.fetchYarnDeps {
              yarnLock = ./yarn.lock;
              hash = "sha256-TkwcsqFLycTtcl+IqPZf5vUQKfelXmJhe1uq+US16lI=";
            };

            buildPhase = ''
              runHook preBuild
              export NG_CLI_ANALYTICS=false
              yarn build
              runHook postBuild
            '';

            installPhase = ''
              runHook preInstall
              mkdir -p "$out"
              cp -r dist/ui/. "$out"
              runHook postInstall
            '';
          };

          docker = pkgs.dockerTools.buildLayeredImage {
            name = "go-pkg-repository-ui";
            tag = dockerTag;

            contents = [ pkgs.dockerTools.binSh pkgs.coreutils pkgs.nodejs_26 entrypoint ];

            extraCommands = ''
              mkdir -p app
              cp -r ${app}/. app/
              chmod 0777 app/browser

              mkdir -p etc
              echo 'appuser:x:10001:10001:Application user:/nonexistent:/sbin/nologin' > etc/passwd
              echo 'appuser:x:10001:' > etc/group
            '';

            config = {
              Entrypoint = [ "${entrypoint}/bin/go-pkg-repository-ui-entrypoint" ];
              User = "10001:10001";
              WorkingDir = "/app";
              ExposedPorts."4000/tcp" = { };
            };
          };
        in
        {
          inherit app docker;
          default = docker;
        });

      devShells = forAllDevSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            packages = [
              pkgs.nodejs_26
              (pkgs.yarn.override { nodejs = pkgs.nodejs_26; })
              (pkgs.writeShellScriptBin "ng" ''
                exec "$PWD/node_modules/.bin/ng" "$@"
              '')
            ];
          };
        });
    };
}
