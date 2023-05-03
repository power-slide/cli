#!/usr/bin/env bash
changelog="CHANGELOG.md"
echo -n "Generating $changelog... "
echo "# Change Log" > $changelog

git fetch --tags
releases=$(git tag -l --sort=-v:refname)
while IFS= read -r tag || [[ -n $tag ]]; do
    gh release view $tag --json body --jq .body >> $changelog
done < <(printf '%s' "$releases")

echo "Done!"
