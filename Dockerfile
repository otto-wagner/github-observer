FROM golang:1.21.9-alpine AS build

WORKDIR /app

COPY . .

RUN go build -a -o observer cmd/main.go

FROM alpine:3.19.1

ARG USERNAME=observer
ARG USER_UID=1000
ARG GITHUB_TOKEN
ENV GITHUB_TOKEN=$GITHUB_TOKEN

RUN touch executor.log && chown $USER_UID:$USER_UID executor.log
RUN touch watcher.log && chown $USER_UID:$USER_UID watcher.log
RUN touch listener.log && chown $USER_UID:$USER_UID listener.log

RUN adduser -u $USER_UID -D $USERNAME
USER $USERNAME

COPY --from=build /app/observer /observer
COPY --from=build /app/conf /conf

ENTRYPOINT ["/observer", "server", "run"]