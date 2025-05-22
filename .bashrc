# export $(grep -v '^#' taskfile.env | xargs)
# alias terminalizer='/home/user/node_modules/.bin/terminalizer'
# alias r='terminalizer record -k'
# alias p='terminalizer play'
# s() {
#     local arg1="$1"

#     terminalizer render "${arg1}" --output "/home/user/gifs/${arg1}"
# }

# # source ~/.bashrc


# . /home/user/.bash_env

alias r='termtosvg'
alias rr='vhs /home/user/tapes/preview.tape'

export PATH=$PATH:/home/user/.local/bin

function all() {
  local tapes_dir="./tapes"

  # Check if the tapes directory exists
  if [ ! -d "$tapes_dir" ]; then
    echo "Error: Directory '$tapes_dir' not found."
    return 1
  fi

  # Count how many tape files we found
  local tape_count=$(find "$tapes_dir" -name "*.tape" | wc -l)

  if [ "$tape_count" -eq 0 ]; then
    echo "No *.tape files found in $tapes_dir"
    return 1
  fi

  echo "Found $tape_count tape files. Processing..."

  # Process each tape file
  find "$tapes_dir" -name "*.tape" | while read tape_file; do
    echo "→ Processing: $tape_file"
    vhs "$tape_file"

    # Check if vhs command was successful
    if [ $? -eq 0 ]; then
      echo "✓ Successfully processed: $tape_file"
    else
      echo "✗ Failed to process: $tape_file"
    fi
    echo "----------------------------------"
  done

  echo "All tape files processed!"
}

function pall() {
  local tapes_dir="./tapes"

  if [ ! -d "$tapes_dir" ]; then
    echo "Error: Directory '$tapes_dir' not found."
    return 1
  fi

  local tape_files
  tape_files=$(find "$tapes_dir" -name "*.tape")
  local tape_count=$(echo "$tape_files" | wc -l)

  if [ "$tape_count" -eq 0 ]; then
    echo "No *.tape files found in $tapes_dir"
    return 1
  fi

  echo "Found $tape_count tape files. Processing in parallel..."

  echo "$tape_files" | xargs -P "$(nproc)" -I {} bash -c '
    echo "→ Processing: {}"
    if vhs "{}"; then
      echo "✓ Successfully processed: {}"
    else
      echo "✗ Failed to process: {}"
    fi
    echo "----------------------------------"
  '

  echo "All tape files processed!"
}
