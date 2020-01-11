############################
# STEP 1 build executable binary
#https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
############################
FROM golang:alpine AS builder


#http://smartystreets.com/blog/2018/09/private-dependencies-in-docker-and-go
#then docker build --build-arg DOCKER_GIT_CREDENTIALS -t makako-gateway .
#ARG DOCKER_GIT_CREDENTIALS

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

#RUN git config --global credential.helper store && echo "${DOCKER_GIT_CREDENTIALS}" > ~/.git-credentials
#RUN git config --global url."https://jebo87:REPB8bsG7TWPuBHzyS9n@bitbucket.org/".insteadOf "https://bitbucket.org/"
WORKDIR $GOPATH/src/gitlab.com/jebo87/makako-gateway/
COPY . .

# Fetch dependencies.
# Using go get.
RUN go get -d -v 
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /makako-gateway/bin/makako-gateway

############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /makako-gateway/bin/makako-gateway /makako-gateway/bin/makako-gateway

# Run the hello binary.
ENTRYPOINT ["/makako-gateway/bin/makako-gateway", "-deployed=true"]
EXPOSE 8087
#export DOCKER_GIT_CREDENTIALS="$(cat ~/.git-credentials)"
#docker build --build-arg DOCKER_GIT_CREDENTIALS -t makako-gateway:0.1 .
#docker run --rm --name makako-gateway --network makako-network -v $(pwd)/config:/makako-gateway/bin/config -p 8087:8087 makako-gateway:0.1
#docker run -d --name makako-gateway --network makako-network -v $(pwd)/config:/makako-gateway/bin/config -p 8087:8087 makako-gateway:0.1