{ ... }:
{
  projectRootFile = "flake.nix";
  programs.gofmt.enable = true;
  programs.nixfmt.enable = true;
  programs.deadnix.enable = true;
}