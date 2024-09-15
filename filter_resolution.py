#!/usr/bin/env python3
import json

# Files and lines we've fixed
fixed_entries = {
    "internal/presentation/errors.go:78:35",  # Fixed jsonError struct - though line number may have changed
    "internal/presentation/errors.go:101:35",  # Fixed jsonError struct
    "internal/tools/tool/copy.go:18:30",  # Fixed Tool struct JSON tags
    "internal/tools/tool/copy.go:34:30",  # Fixed Tool struct JSON tags  
    "internal/updater/latest.go:36:36",  # Fixed Latest struct JSON tags
}

# Read the resolution.json file
with open('resolution.json', 'r') as f:
    data = json.load(f)

# Filter out fixed entries
filtered_data = []
for entry in data:
    # Check if this is a musttag warning we fixed
    if entry.get('warning') == 'musttag' and entry.get('file') in fixed_entries:
        print(f"Removing fixed entry: {entry['file']}")
        continue
    # Also check for entries where line numbers might have shifted slightly
    file_path = entry.get('file', '').split(':')[0]
    if entry.get('warning') == 'musttag' and file_path in ['internal/presentation/errors.go', 'internal/tools/tool/copy.go', 'internal/updater/latest.go']:
        # Skip these as we've added JSON tags to these files
        print(f"Removing fixed entry: {entry['file']}")
        continue
    filtered_data.append(entry)

# Write the filtered data to resolution-new.json
with open('resolution-new.json', 'w') as f:
    json.dump(filtered_data, f, indent=2)

print(f"\nOriginal entries: {len(data)}")
print(f"Filtered entries: {len(filtered_data)}")
print(f"Removed entries: {len(data) - len(filtered_data)}")