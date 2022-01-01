#!/bin/bash

 go build -ldflags="-X 'main.BuildTime=$(date)'"