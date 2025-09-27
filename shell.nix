# Shell enviorment for nix package manager
let
  pkgs = import <nixpkgs> { config.allowUnfree = true; };
in pkgs.mkShell {
  packages = [
    pkgs.turso-cli
    pkgs.google-cloud-sdk-gce
  ];
  shellHook = ''
    export PS1="\n\[\033[1;32m\][ci-cd-learning:\w]\$\[\033[0m\]"
  '';
}
