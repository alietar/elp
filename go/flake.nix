{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-25.11";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils, lib, ... }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };

        myGoApp = pkgs.buildGoModule {
          pname = "go-reachable";
          version = "0.0.1";
          src = ./.;

          shellHook = ''
            export GOFLAGS="-mod=readonly"
          '';

          vendorHash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
        };
      in
      {
        packages.default = myGoApp;

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ go gopls gotools go-tools _7zz graphviz ];
        };
      }) // {
        nixosModules.default = { config, lib, pkgs, ... }:
          let
            cfg = config.services.go-reachable;
          in
          {
            options.services.go-reachable = {
              enable = lib.mkEnableOption "Activer le service Go";
              extraArgs = lib.mkOption {
                type = lib.types.listOf lib.types.str;
                default = [];
                description = "Arguments supplémentaires à passer au binaire.";
              };
            };

            config = lib.mkIf cfg.enable {
              systemd.services.go-reachable = {
                description = "Mon Service Go";
                after = [ "network.target" ];
                wantedBy = [ "multi-user.target" ];
                serviceConfig = {
                  PreExec = "${self.packages.${pkgs.system}.default}/bin/go-reachable -accuracy-5 -dl-some 69";
                  ExecStart = "${self.packages.${pkgs.system}.default}/bin/go-reachable ${lib.escapeShellArgs cfg.extraArgs}";
                  Restart = "always";
                  User = "nobody";
                };
              };
            };
          };
      };
}
