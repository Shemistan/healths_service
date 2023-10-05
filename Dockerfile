ARG BUILDER_IMAGE=golang:1.21.1-alpine
ARG ALPINE_IMAGE=alpine:3.15

FROM $BUILDER_IMAGE AS build

COPY . /app/
WORKDIR /app

RUN go mod download
RUN go build -o /bin/healths_service ./cmd/main.go

FROM $ALPINE_IMAGE

WORKDIR /root/
COPY --from=build /bin/healths_service .
COPY --from=build /app/dev/local.env dev/local.env
COPY --from=build /app/config/api_list.txt config/api_list.txt

CMD ["./healths_service"]