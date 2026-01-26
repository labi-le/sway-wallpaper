{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    golangci-lint
    goreleaser
    gopls
  ];

  shellHook = ''
    echo "Welcome to the Go development environment for chiasma"
    go version
  '';
}
