# BUILD
FROM golang:1.16-alpine as build

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .


# RUN
FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /app/main .

EXPOSE 8080
CMD [ "./main" ]