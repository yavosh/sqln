FROM golang:1.17 as builder

WORKDIR /build
COPY ./ ./

# ldflags -w and -s result in smaller binary.
# -w Omits the DWARF symbol table.
# -s Omits the symbol table and debug information.
RUN CGO_ENABLED=1 go build -mod=mod -a -ldflags "-w -s" -o sqln ./cmd/sqln/main.go
RUN CGO_ENABLED=0 go build -mod=mod -a -ldflags "-w -s" -o sqln-cli ./cmd/sqln-cli/main.go

FROM scratch
COPY --from=builder /build/sqln /
COPY --from=builder /build/sqln-cli /

ADD ./build/ca-certificates.crt /etc/ssl/certs/
ADD ./build/zoneinfo.zip /zoneinfo.zip

ENV ZONEINFO "/zoneinfo.zip"

CMD ["/sqln"]