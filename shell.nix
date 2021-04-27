{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [ gcc go gopls golangci-lint ];
  CGO_ENABLED = "1";
}
