FROM golang:alpine3.14 as build
# set env
ARG github_key
ENV github_key=${github_key}

# github preparation
RUN apk --no-cache add git;apk --no-cache add ca-certificates
RUN git config --global url."${github_key}".insteadOf "https://github.com"
# go env
RUN go env -w GOPRIVATE=github.com/mindtera

RUN mkdir -p /app
WORKDIR /app

# Go Modules preparation
COPY . .
RUN cd cmd && go build -o ../main-app

FROM alpine
ENV TZ=GMT
RUN apk --no-cache add ca-certificates
COPY --from=build /app/main-app /main-app
COPY --from=build /app/secret /secret

ENTRYPOINT ["/main-app"]