FROM golang:1.21.9-alpine AS build

ARG ACTIVATE_SSL=false

ARG COUNTRY=DE
ARG STATE=Hamburg
ARG COMMON_NAME=localhost
ARG ORGANIZATION=OpenSourceOrg
ARG ORG_UNIT=OpenSourceUnit

ENV COUNTRY=$COUNTRY \
    STATE=$STATE \
    LOCALITY=$LOCALITY \
    COMMON_NAME=$COMMON_NAME \
    ORGANIZATION=$ORGANIZATION \
    ORG_UNIT=$ORG_UNIT

WORKDIR /app

COPY . .

RUN if [ "$ACTIVATE_SSL" = "true" ]; then \
     chmod +x ./scripts/generate_certificate.sh; \
     apk add -U --no-cache ca-certificates git openssl; \
     ./scripts/generate_certificate.sh; \
    fi

RUN go build -a -o observer cmd/main.go

FROM alpine:3.19.1

ARG USERNAME=observer
ARG USER_UID=1000
ARG GITHUB_TOKEN
ENV GITHUB_TOKEN=$GITHUB_TOKEN

RUN touch executor.log && chown $USER_UID:$USER_UID executor.log
RUN touch watcher.log && chown $USER_UID:$USER_UID watcher.log
RUN touch listener.log && chown $USER_UID:$USER_UID listener.log

COPY --from=build /app/certs /certs
RUN chown -R $USER_UID:$USER_UID /certs
COPY --from=build /app/conf /conf
RUN chown -R $USER_UID:$USER_UID /conf

RUN adduser -u $USER_UID -D $USERNAME
USER $USERNAME

COPY --from=build /app/observer /observer

ENTRYPOINT ["/observer", "server", "run"]