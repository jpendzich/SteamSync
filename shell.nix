{ pkgs ? import <nixpkgs> {} }:
    pkgs.mkShell {
        nativeBuildInputs = with pkgs.buildPackages; [ 
            go
            gopls
            typescript-language-server
            vscode-langservers-extracted
        ];
}
