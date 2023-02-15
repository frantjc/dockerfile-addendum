ARG base_image=scratch
ARG build_image=golang:1.20-alpine3.16

FROM ${base_image} AS base_image

FROM ${build_image} AS build_image
ENV CGO_ENABLED=0

FROM build_image AS build
WORKDIR $GOPATH/src/github.com/frantjc/dockerfile-addendum
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG version=0.0.0
ARG prerelease=
RUN go build -ldflags "-s -w -X github.com/frantjc/dockerfile-addendum.Version=${version} -X github.com/frantjc/dockerfile-addendum.Prerelease=${prerelease}" -o /usr/local/bin ./cmd/addendum

FROM base_image AS addendum
COPY --from=build /usr/local/bin /
ENTRYPOINT ["addendum"]

FROM addendum
