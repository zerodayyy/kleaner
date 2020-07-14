FROM golang:1.14
WORKDIR /go/src/github.com/zerodayyy/kleaner/
COPY go.mod ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o kleaner .

FROM scratch
WORKDIR /
COPY --from=0 /go/src/github.com/zerodayyy/kleaner/kleaner ./
ENTRYPOINT ["/kleaner"]
