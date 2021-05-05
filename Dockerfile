############################
# STEP 1 build executable binary
#https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
############################
FROM golang:alpine AS builder

ENV GO111MODULE=on

#Receive the git credentials as an argunment
ARG DOCKER_GIT_CREDENTIALS

# Install git.
# Git is required for fetching the dependencies.
#RUN apk update && apk add --no-cache git
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

#set the git config
RUN git config --global credential.helper store && echo "${DOCKER_GIT_CREDENTIALS}" > ~/.git-credentials
RUN git config --global url."${DOCKER_GIT_CREDENTIALS}/".insteadOf "https://gitlab.com/"
WORKDIR $GOPATH/src/gitlab.com/jebo87/makako-gateway/
COPY . .

# Fetch dependencies.
RUN go mod download
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /makako-gateway/bin/makako-gateway

############################
# STEP 2 build a small image
############################
FROM alpine:3.7
# Copy our static executable.
COPY --from=builder /makako-gateway/bin/makako-gateway /makako-gateway/bin/makako-gateway
#pass any environment variables needed by the app
ENV ELASTIC_ADDRESS = $ELASTIC_ADDRESS_ENV
# Run the binary.
ENTRYPOINT ["/makako-gateway/bin/makako-gateway", "-deployed=true"]
EXPOSE 8087
