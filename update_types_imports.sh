#!/bin/bash

# Script to update all type imports to use the new consolidated types structure

echo "Updating type imports to use new consolidated types structure..."

# Update facts types imports
find . -name "*.go" -exec sed -i 's|spookyfactstypes "spooky/internal/facts/types"|spookytypes "spooky/internal/types"|g' {} \;
find . -name "*.go" -exec sed -i 's|"spooky/internal/facts/types"|"spooky/internal/types/facts"|g' {} \;

# Update config types imports
find . -name "*.go" -exec sed -i 's|spookyconfigtypes "spooky/internal/config/types"|spookytypes "spooky/internal/types"|g' {} \;
find . -name "*.go" -exec sed -i 's|"spooky/internal/config/types"|"spooky/internal/types/config"|g' {} \;

# Update actions types imports
find . -name "*.go" -exec sed -i 's|spookyactionstypes "spooky/internal/actions/types"|spookytypes "spooky/internal/types"|g' {} \;
find . -name "*.go" -exec sed -i 's|"spooky/internal/actions/types"|"spooky/internal/types/actions"|g' {} \;

# Update machines types imports
find . -name "*.go" -exec sed -i 's|spookymachinestypes "spooky/internal/machines/types"|spookytypes "spooky/internal/types"|g' {} \;
find . -name "*.go" -exec sed -i 's|"spooky/internal/machines/types"|"spooky/internal/types/machines"|g' {} \;

# Update ssh types imports
find . -name "*.go" -exec sed -i 's|spookysshtypes "spooky/internal/ssh/types"|spookytypes "spooky/internal/types"|g' {} \;
find . -name "*.go" -exec sed -i 's|"spooky/internal/ssh/types"|"spooky/internal/types/ssh"|g' {} \;

# Update templates types imports
find . -name "*.go" -exec sed -i 's|spookytemplatestypes "spooky/internal/templates/types"|spookytypes "spooky/internal/types"|g' {} \;
find . -name "*.go" -exec sed -i 's|"spooky/internal/templates/types"|"spooky/internal/types/templates"|g' {} \;

# Update variables types imports
find . -name "*.go" -exec sed -i 's|spookyvariablestypes "spooky/internal/variables/types"|spookytypes "spooky/internal/types"|g' {} \;
find . -name "*.go" -exec sed -i 's|"spooky/internal/variables/types"|"spooky/internal/types/variables"|g' {} \;

# Update logging types imports
find . -name "*.go" -exec sed -i 's|spookyloggingtypes "spooky/internal/logging/types"|spookytypes "spooky/internal/types"|g' {} \;
find . -name "*.go" -exec sed -i 's|"spooky/internal/logging/types"|"spooky/internal/types/logging"|g' {} \;

# Update secrets types imports
find . -name "*.go" -exec sed -i 's|spookysecretstypes "spooky/internal/secrets/types"|spookytypes "spooky/internal/types"|g' {} \;
find . -name "*.go" -exec sed -i 's|"spooky/internal/secrets/types"|"spooky/internal/types/secrets"|g' {} \;

# Update schemas types imports
find . -name "*.go" -exec sed -i 's|"spooky/internal/schemas/types"|"spooky/internal/types/schemas"|g' {} \;

# Update cli types imports
find . -name "*.go" -exec sed -i 's|"spooky/internal/cli/types"|"spooky/internal/types/cli"|g' {} \;

echo "Type imports updated successfully!"
