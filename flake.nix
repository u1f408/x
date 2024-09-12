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
        pkgs = import inputs.nixpkgs { inherit system; };
        inherit (pkgs) lib;

      in {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
          ];
        };
      });
}
