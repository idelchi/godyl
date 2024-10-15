#!/bin/bash

for file in ~/.local/bin/*; do
  if [[ -x "$file" && ! -d "$file" ]]; then
    echo "Running $file"
    "$file" --version 2>/dev/null
  fi
done
