FROM golang:alpine
RUN apk add --no-cache nodejs npm
WORKDIR /go/src/lyrid-sd/
COPY . .
ENV GO111MODULE=on
RUN go mod tidy
RUN CGO_ENABLED=0 go build -v -o app

WORKDIR /go/src/lyrid-sd/web
RUN npm install
RUN npm run build

FROM alpine
RUN apk add --no-cache ca-certificates bash
WORKDIR /lyrid-sd/
COPY --from=0 /go/src/lyrid-sd/app .
COPY --from=0 /go/src/lyrid-sd/.env .
COPY --from=0 /go/src/lyrid-sd/web/build ./web/build
ENTRYPOINT ["/lyrid-sd/app"]