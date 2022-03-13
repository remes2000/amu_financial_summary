#!/bin/bash

cd frontend/amu-financial-summary
ng build --configuration production --base-href=/app/ --output-path=../dist
cd ../../
env GOOS=freebsd GOARCH=amd64 go install
rm -rf dist
mkdir -p dist
cp ../../../../bin/freebsd_amd64/amu_financial_summary ./dist/
mkdir -p dist/frontend
cp -r ./frontend/dist/** ./dist/frontend