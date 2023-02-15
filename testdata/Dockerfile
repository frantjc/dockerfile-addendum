FROM alpine:3.16
COPY --from=ghcr.io/frantjc/dockerfile-addendum /addendum /usr/local/bin
ARG tgz=test.tgz
ADD ${tgz} /usr/local/bin/
RUN addendum -ru /usr/local/bin/$(basename ${tgz})
