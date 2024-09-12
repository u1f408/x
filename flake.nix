{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/release-24.05";
    fu.url = "github:numtide/flake-utils";
    compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = { self, fu, ... }@inputs:
    fu.lib.eachDefaultSystem (system:
      let
        version = "20240912";

        pkgs = import inputs.nixpkgs { inherit system; };
        inherit (pkgs) lib;

      in {
        packages.u1f408-x = pkgs.buildGoModule {
          pname = "u1f408-x";
          inherit version;

          src = ./.;
          vendorHash = "sha256-kSP70NpnEOyeQjsuNjqOdsvdwQtebLl5oH1ySadxsUI=";

          ldflags = [ "-X main.Version=${version}" ];
          subPackages = [
            "proxyssh"
            "box"
          ];
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gnumake
          ];
        };
      });
}
