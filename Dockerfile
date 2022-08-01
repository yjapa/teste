FROM debian:buster

RUN apt-get update \
    && apt-get install -y ca-certificates --no-install-recommends \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && update-ca-certificates 

RUN	addgroup klever \
    && adduser --disabled-password --gecos "" klever --uid 1000 --ingroup klever \
    && chown klever:root /etc/ssl/private

USER klever

COPY --chown=klever:klever app /usr/local/bin/

RUN chmod +x /usr/local/bin/app
ENTRYPOINT [ "/usr/local/bin/app" ]