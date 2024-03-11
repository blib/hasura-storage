#!/usr/bin/env bash

which nix >/dev/null

IMAGE=${1:-"dockerImage"}
SYSTEM=${2:-""}

if [[ $NIX_BUILD_NATIVE -eq 1 ]]; then
    if [[ -z $SYSTEM ]]; then
        case $(uname -m) in
        "arm64")
            SYSTEM="aarch64-linux"
            ;;
        *)
            SYSTEM="x86_64-linux"
            ;;
        esac
    fi

    nix --extra-experimental-features nix-command --extra-experimental-features flakes build .\#packages.${SYSTEM}.$IMAGE --print-build-logs && docker load <result
    exit $?
fi

if [[ ($? -eq 0) && ($(uname) == "Linux") ]]; then
    nix --extra-experimental-features nix-command --extra-experimental-features flakes build .\#$IMAGE --print-build-logs && docker load <result
    exit $?
fi

docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v $PWD:/build \
    -w /build \
    --entrypoint sh \
    dbarroso/nix:2.6.0 \
    -c "nix build .\\#packages.$IMAGE --print-build-logs && docker load < result"
