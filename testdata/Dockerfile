FROM scratch
ARG tgz=zip_3.0_x86_64.tgz
ADD ${tgz} /usr/local/bin/test.tgz
COPY --from=ghcr.io/frantjc/dockerfile-addendum /addendum /usr/local/bin
RUN addendum -ruo /usr/local/bin /usr/local/bin/test.tgz
