set -e

staged=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -z "$staged" ]; then
  echo "No Go files changed, skipping lint."
  exit 0
fi

for pkg in $(echo "$staged" | xargs -n1 dirname | sort -u); do
  echo "→ Linting package: $pkg"
  golangci-lint run "$pkg"
done