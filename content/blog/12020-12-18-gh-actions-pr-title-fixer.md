---
description: spinning up a vm just to run a bash script...
title: github action pr title fixer
---

### _github_ actions

Say you need your pull request titles to fit a certain format
(prefix with `WD-$ticketnumber`) because your automation relies on it.
Also you want a clickable link in the body.
You could write an action to fix them for you:

todo:

- logic to determine if a patch isn't actually necessary

```yaml
name: PR Title Fixer

on:
  pull_request:
    types:
      - opened

jobs:
  fix:
    name: Fix
    runs-on: ubuntu-latest
    steps:
      # only if we can use the branch as the source of truth
      - if: startsWith(${{ github.head_ref }}, "WD-")
        run: |
          # inline script
          set -euxo pipefail

          branch="${{ github.head_ref }}"

          # extract WD-ticket from potentially long name: WD-ticket-some-description
          if [[ "$branch" =~ ^(WD-[0-9]{5}).* ]] ; then
            branch=${BASH_REMATCH[1]}
            echo "fixed branch is $branch"
          else
            # not in expected format
            exit 0
          fi

          cat << EOF | \
            jq \
              --arg b $branch \
              --arg u "jira: [$branch](https://example.com/browse/$branch)" \
             '{title: [$b, (.title | sub("^[Ww][Dd][- ][0-9]{5}"; ""))] | join(" "), body: [$u, .body] | join("\n\n")}' | \
            curl \
              --silent \
              --request PATCH \
              --header "Content-Type: application/json"  \
              --header "Authorization: Bearer ${{ github.token }}" \
              --data-binary @- \
              "https://api.github.com/repos/${{ github.repository }}/pulls/${{ github.event.number }}" > /dev/null
          ${{ toJson(github.event.pull_request) }}
          EOF
```
