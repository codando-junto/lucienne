#!/bin/bash
sleep 10
npm install
node esbuild.js &
/go/bin/air
