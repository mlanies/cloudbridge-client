#!/bin/bash

# Script to bump version and create a release

set -e

# Check if version type is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <major|minor|patch>"
    echo "Example: $0 patch"
    exit 1
fi

VERSION_TYPE=$1
CURRENT_VERSION=$(cat VERSION)

# Parse current version
IFS='.' read -ra VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR=${VERSION_PARTS[0]}
MINOR=${VERSION_PARTS[1]}
PATCH=${VERSION_PARTS[2]}

# Bump version based on type
case $VERSION_TYPE in
    major)
        NEW_MAJOR=$((MAJOR + 1))
        NEW_VERSION="${NEW_MAJOR}.0.0"
        ;;
    minor)
        NEW_MINOR=$((MINOR + 1))
        NEW_VERSION="${MAJOR}.${NEW_MINOR}.0"
        ;;
    patch)
        NEW_PATCH=$((PATCH + 1))
        NEW_VERSION="${MAJOR}.${MINOR}.${NEW_PATCH}"
        ;;
    *)
        echo "Invalid version type. Use: major, minor, or patch"
        exit 1
        ;;
esac

echo "Current version: $CURRENT_VERSION"
echo "New version: $NEW_VERSION"

# Update VERSION file
echo "$NEW_VERSION" > VERSION

# Commit the version bump
git add VERSION
git commit -m "Bump version to $NEW_VERSION"

# Create and push tag
git tag -a "v$NEW_VERSION" -m "Release v$NEW_VERSION"
git push origin main
git push origin "v$NEW_VERSION"

echo "Version bumped to $NEW_VERSION and tag v$NEW_VERSION created"
echo "GitHub Actions will now create a release automatically" 