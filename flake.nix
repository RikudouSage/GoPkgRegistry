{
  description = "Go package repository server and container image";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      versionFromEnv = builtins.getEnv "PACKAGE_VERSION";
      packageVersion = if versionFromEnv == "" then "0.0.0" else versionFromEnv;
      dockerTag = if versionFromEnv == "" then "dev" else versionFromEnv;
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};

          buildGo127Module = pkgs.buildGoModule.override {
            go = pkgs.go_1_27;
          };

          app = buildGo127Module {
            pname = "go-pkg-repository";
            version = packageVersion;
            src = self;

            env.CGO_ENABLED = "1";
            ldflags = [ "-s" "-w" ];
            vendorHash = "sha256-6IsWj4ANSbqSSpYPTeAqqYfcoEb8DURJh69l6NRWD8k=";
          };

          docker = pkgs.dockerTools.buildLayeredImage {
            name = "go-pkg-repository";
            tag = dockerTag;

            contents = [
              pkgs.cacert
              pkgs.tzdata
            ];

            extraCommands = ''
              mkdir -p etc
              cp ${app}/bin/go-pkg-repository app
              chmod 0555 app

              echo 'appuser:x:10001:10001:Application user:/nonexistent:/sbin/nologin' > etc/passwd
              echo 'appuser:x:10001:' > etc/group
            '';

            config = {
              Entrypoint = [ "/app" ];
              User = "10001:10001";
              ExposedPorts."8080/tcp" = { };
            };
          };
        in
        {
          inherit app docker;
          default = docker;
        });
    };
}
