FROM golang:1.15 as build
RUN mkdir /ethminer-exporter
WORKDIR /ethminer-exporter
COPY . ./
RUN GOOS=linux CGO_ENABLED=0 GO111MODULE=on go build -mod=vendor -a -tags netgo -ldflags '-w' -o /bin/ethxporter main.go

FROM scratch as final
COPY --from=build /bin/ethxporter /
ENTRYPOINT [ "/ethxporter" ]
CMD ["--help"]
