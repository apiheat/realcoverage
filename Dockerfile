FROM amd64/alpine:latest

ARG binary_version

ENV AKAMAI_CLI_HOME=/cli

RUN apk add --no-cache --update openssl \
        ca-certificates \
        libc6-compat \
        libstdc++ \
        wget \
        curl \
        jq \
        bash \
        nodejs \
        npm \
        rm -rf /var/cache/apk/* && \
    wget --quiet -O /usr/local/bin/akamai https://github.com/akamai/cli/releases/download/1.1.4/akamai-1.1.4-linuxamd64 && \
    chmod +x /usr/local/bin/akamai && \
    echo '[ ! -z "$TERM" -a -r /etc/motd ] && cat /etc/motd' >> /root/.bashrc


RUN mkdir -p /cli/.akamai-cli && \
    echo "[cli]" > /cli/.akamai-cli/config && \
    echo "cache-path            = /cli/.akamai-cli/cache" >> /cli/.akamai-cli/config && \
    echo "config-version        = 1" >> /cli/.akamai-cli/config && \
    echo "enable-cli-statistics = false" >> /cli/.akamai-cli/config && \
    echo "last-ping             = 2018-04-27T18:16:12Z" >> /cli/.akamai-cli/config && \
    echo "client-id             =" >> /cli/.akamai-cli/config && \
    echo "install-in-path       =" >> /cli/.akamai-cli/config && \
    echo "last-upgrade-check    = ignore" >> /cli/.akamai-cli/config

RUN akamai install property --force && \
    rm -rf /cli/.akamai-cli/src/akamai-cli-netlist/.git
RUN akamai install https://github.com/apiheat/akamai-cli-overview --force && \
    rm -rf /cli/.akamai-cli/src/akamai-cli-overview/.git
# RUN wget --quiet -O /usr/local/bin/realcoverage https://github.com/apiheat/realcoverage/releases/download/v$binary_version/realcoverage_linux_amd64 && \
#    chmod +x /usr/local/bin/realcoverage

ENV AKAMAI_CLI_HOME=/cli
VOLUME /cli
VOLUME /root/.edgerc

CMD ["/bin/bash"]