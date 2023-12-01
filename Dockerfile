FROM golang:1.18

COPY . /work
WORKDIR /work
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nginxlog-exporter

FROM scratch

COPY --from=0 /work/nginxlog-exporter /nginxlog-exporter

EXPOSE 7788
ENTRYPOINT ["/nginxlog-exporter"]
