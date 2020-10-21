#!/usr/bin/env echo The following file must be sourced to affect $PATH:

# Assumes all the cmd packages are `go install`ed to the same directory.
export PATH="$PATH:$(dirname $(go list -f '{{.Target}}' nat_project/cmd/... | head -n 1))"
