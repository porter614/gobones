### Build Stage ###
FROM golang:latest AS build
ARG APP

WORKDIR /${APP}

# First get the go.mod and download dependencies
COPY go.mod ./
RUN go mod download

# Grab source code...
COPY . /${APP}

# Build it!
RUN go build -mod=readonly

### Lint Stage ###
FROM build AS lint

RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.16.0

### Unit Test Stage ###
FROM lint AS test
ARG BUILD_NUMBER
LABEL build=${BUILD_NUMBER} \
   image=coverage

RUN go get golang.org/x/tools/cmd/cover && \
   go get github.com/t-yuki/gocover-cobertura && \
   go test -coverprofile=coverage.txt -covermode count ./pkg/... && \
   $GOPATH/bin/gocover-cobertura < coverage.txt > coverage.xml

### Cross-compile Stage ###
FROM build AS cross-compile
ARG IMPORT_PATH
ARG VERSION
ARG APP
ARG SEMVER
ARG BUILD_DATE
ARG BUILD_USER
ARG GIT_COMMIT
ARG GIT_BRANCH

# Cross compile for scratch
RUN CGO_ENABLED=0 GOOS=linux go build \
   -a -installsuffix cgo \
   -ldflags "-X ${IMPORT_PATH}/pkg/version.App=${APP} \
             -X ${IMPORT_PATH}/pkg/version.Version=${VERSION} \
             -X ${IMPORT_PATH}/pkg/version.SemVer=${SEMVER} \
             -X ${IMPORT_PATH}/pkg/version.BuildDate=${BUILD_DATE} \
             -X ${IMPORT_PATH}/pkg/version.BuildUser=${BUILD_USER} \
             -X ${IMPORT_PATH}/pkg/version.CommitId=${GIT_COMMIT} \
             -X ${IMPORT_PATH}/pkg/version.Branch=${GIT_BRANCH} \
             -X ${IMPORT_PATH}/cmd.App=${APP}" \
   -o /go/bin/app

# Create the porter614 username
RUN useradd -ms /bin/bash porter614
RUN adduser porter614 porter614

### Final Stage ###
FROM scratch AS deploy

# Import the user and group files from the builder.
COPY --from=cross-compile /etc/passwd /etc/passwd
COPY --from=cross-compile /etc/group /etc/group

# We need the app
COPY --from=cross-compile /go/bin/app /bin/

# We need the config
ADD --chown=porter614:porter614 config/config.json /config/

# Use an unprivileged user
USER porter614

# Run it
ENTRYPOINT ["/bin/app"]
CMD ["run"]
