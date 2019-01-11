FROM quay.io/prometherion/dep:0.5.0 AS build
WORKDIR /go/src/github.com/prometherion/prometheus-exporter-golang/app
ADD ./app ./
RUN dep ensure -v && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /tmp/prometheus-exporter-golang .

FROM scratch
COPY --from=build /tmp/prometheus-exporter-golang /prometheus-exporter-golang
CMD ["/prometheus-exporter-golang"]