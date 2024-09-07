{ lib, buildGoModule }:

buildGoModule rec {
  pname = "nes-chr";
  version = "0.0.1";

  src = ./.;

  vendorHash = null;

  meta = with lib; {
    description = "Extract the NES CHR ROM graphics";
    homepage = "https://github.com/theobori/${pname}";
    license = licenses.mit;
  };
}