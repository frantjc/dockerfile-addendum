# dockerfile-ADDendum

> _not yet thoroughly tested_

Small script to pair with the Dockerfile `ADD` directive to make an idempotent action regardless of if the `<src>` of the `ADD` is a local tar archive or a tar archive from a remove URL.

## use case

A Dockerfile's `ADD` directive is a bit [finicky](https://docs.docker.com/engine/reference/builder/#add). It allows files to be added to the image either from a local directory or a remote URL. This has value in that if the `<src>` of the `ADD` is an `ARG`, then the functionality of the `ADD` can be changed with build-time arguments.

For example, this allows developers to lazily get the latest version of the dependency from a tar archive at some remote URL while CI consistently supplies the same tar archive to the Dockerfile as a `--build-arg`.

However, if that `<src>` is expected to be either a local tar archive or a remote tar archive, the `ADD` directive is not idempotent; `ADD` unpacks local tar archives but not remote ones. As a result, subsequent Dockerfile directives would have to account for the difference to make the build successful.

This is where the ADDendum comes in:

## usage

```Dockerfile
ARG zip=zip_3.0_x86_64.tgz
ADD ${zip} /usr/local/bin
COPY --from=ghcr.io/frantjc/dockerfile-addendum /addendum /usr/local/bin
RUN addendum -ruo /usr/local/bin /usr/local/bin/zip_3.0_x86_64.tgz
```

Now this Dockerfile can be built with `--build-arg zip=<src>` where `<src>` is either a tar archive at a remote URL or a local tar archive.
