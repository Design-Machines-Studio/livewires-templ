# Development container for Go + Templ
FROM golang:1.25-alpine

RUN apk add --no-cache gcc musl-dev git

# Install templ CLI. Pinned to the version that produced the committed
# *_templ.go files: a floating @latest rewrote all 39 of them against a newer
# runtime API (templ.ResolveAttributeValue) than the go.mod runtime provides,
# which broke the build. Bumping this pin and the go.mod runtime together is a
# separate change: it regenerates every *_templ.go in the repo.
RUN go install github.com/a-h/templ/cmd/templ@v0.3.977

WORKDIR /app

# Source code is mounted via volume
CMD ["sh"]
