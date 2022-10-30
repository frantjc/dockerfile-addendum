FROM alpine:3.15
ARG tgz=test.tgz
ADD ${tgz} /usr/local/bin/
COPY --from=ghcr.io/frantjc/dockerfile-addendum /addendum /usr/local/bin
RUN addendum -ruo /usr/local/bin /usr/local/bin/$(basename ${tgz})
