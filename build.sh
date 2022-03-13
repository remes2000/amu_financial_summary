#!/bin/bash

cd frontend/amu-financial-summary
ng build --configuration production --base-href=/app/ --output-path=../dist
cd ../../
go install
rm -rf dist
mkdir -p dist
cp ../../../../bin/amu_financial_summary ./dist/
mkdir -p dist/frontend
cp -r ./frontend/dist/** ./dist/frontend