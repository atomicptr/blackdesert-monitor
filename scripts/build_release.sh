#!/usr/bin/env bash

export GOOS=windows

rm -rf build
mkdir -p build

go build -o build/blackdesert-monitor.exe
cp settings.yaml.skel build/settings.yaml
cp README.md build/README.md

zip -j blackdesert-monitor.zip build/*
