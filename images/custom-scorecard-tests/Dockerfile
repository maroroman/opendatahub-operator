# Base image
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

ENV HOME=/opt/custom-scorecard-tests \
    USER_NAME=custom-scorecard-tests \
    USER_UID=1001

RUN echo "${USER_NAME}:x:${USER_UID}:0:${USER_NAME} user:${HOME}:/sbin/nologin" >> /etc/passwd

WORKDIR ${HOME}

ARG BIN=/bin/custom-scorecard-tests
COPY $BIN /custom-scorecard-tests 
COPY ../../internal/tests /bin

RUN go build -o bin/custom-scorecard-tests main.go

ENTRYPOINT ["/custom-scorecard-tests"]

USER ${USER_UID}