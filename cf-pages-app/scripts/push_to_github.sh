#!/usr/bin/env bash
set -euo pipefail

# Usage:
#   GITHUB_TOKEN=ghp_xxx GITHUB_USER=youruser REPO_NAME=cf-pages-chat ./scripts/push_to_github.sh
#
# Requirements:
#   - curl
#   - git repository already initialized with a main branch

REPO_NAME="${REPO_NAME:-cf-pages-chat}"
GITHUB_USER="${GITHUB_USER:-}"
GITHUB_TOKEN="${GITHUB_TOKEN:-}"

if [ -z "$GITHUB_USER" ] || [ -z "$GITHUB_TOKEN" ]; then
  echo "GITHUB_USER and GITHUB_TOKEN must be set" >&2
  exit 1
fi

API=https://api.github.com

echo "Ensuring repo exists: ${GITHUB_USER}/${REPO_NAME}"
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
  -H "Authorization: token ${GITHUB_TOKEN}" \
  "${API}/repos/${GITHUB_USER}/${REPO_NAME}")

if [ "$HTTP_STATUS" = "404" ]; then
  echo "Creating repo ${REPO_NAME} under ${GITHUB_USER}"
  CREATE_STATUS=$(curl -s -o /dev/null -w "%{http_code}" \
    -H "Authorization: token ${GITHUB_TOKEN}" \
    -H 'Content-Type: application/json' \
    -d "{\"name\":\"${REPO_NAME}\",\"private\":false}" \
    "${API}/user/repos")
  if [ "$CREATE_STATUS" -ge 300 ]; then
    echo "Failed to create repo (status $CREATE_STATUS). Check token scopes (repo)." >&2
    exit 1
  fi
else
  echo "Repo already exists (HTTP $HTTP_STATUS)"
fi

REMOTE_URL="https://${GITHUB_TOKEN}@github.com/${GITHUB_USER}/${REPO_NAME}.git"

if git remote | grep -q "^origin$"; then
  git remote set-url origin "$REMOTE_URL"
else
  git remote add origin "$REMOTE_URL"
fi

git push -u origin main
echo "Pushed to https://github.com/${GITHUB_USER}/${REPO_NAME}"
