export GO111MODULE="on"
export GOPATH="$HOME/go"

# Load environment variables from .env file
ENV_FILE="artifacts/.asset.env"
if [ -f "$ENV_FILE" ]; then
  while IFS= read -r line; do
    # comments and white spaces lines
    if [[ $line != \#* ]] && [[ -n "$line" ]] && [[ ! "$line" =~ ^[[:space:]]*$ ]]; then
      export "$line"
    fi
  done < "$ENV_FILE"
fi

# TODO: implement tags instead
go test -v -run TestAsset ./internal/repository