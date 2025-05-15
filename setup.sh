#!/bin/bash

# Define blocks for each category
declare -A blocks
blocks["forms"]="login signup search settings contact"
blocks["navigation"]="navbar sidebar breadcrumb"
blocks["content"]="card table hero comment"
blocks["layouts"]="base auth dashboard"
blocks["partials"]="formfield pagination alert filtertag"

# Base directory for blocks
base_dir="library/blocks"

# Create blocks directory if it doesn't exist
mkdir -p "$base_dir"

# Loop through each category
for category in "${!blocks[@]}"; do
  # Create category subdirectory
  category_dir="$base_dir/$category"
  mkdir -p "$category_dir"

  # Split blocks string into array and loop
  IFS=' ' read -r -a block_array <<< "${blocks[$category]}"
  for block in "${block_array[@]}"; do
    # Create .templ file
    templ_file="$category_dir/$block.templ"
    # Capitalize first letter of block name for templ function
    block_name="$(echo "$block" | awk '{print toupper(substr($0,1,1)) tolower(substr($0,2))}')"
    cat << EOF > "$templ_file"
package $category
templ $block_name() {}
EOF
    echo "Created $templ_file"

    css_file="$category_dir/$block.css"
    touch "$css_file"
    echo "Created $css_file"
  done
done
