FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/shortlink ./cmd/shortlink

FROM scratch
COPY --from=build /out/shortlink /shortlink
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/shortlink"]
