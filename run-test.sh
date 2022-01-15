#!/bin/bash

go test ./... -json -cover | tparse
