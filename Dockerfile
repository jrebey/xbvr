FROM ubuntu:19.04
RUN apt update && apt install -y wget ca-certificates

ARG DRONE_TAG

RUN wget -O /tmp/xbvr.tgz "https://github.com/xbapps/xbvr/releases/download/"$DRONE_TAG"/xbvr_"$DRONE_TAG"_Linux_x86_64.tar.gz" && \
    tar xvfz /tmp/xbvr.tgz -C /usr/local/bin/ && \
    rm /tmp/xbvr.tgz

EXPOSE 9998-9999
VOLUME /root/.config/

ENTRYPOINT ["/usr/local/bin/xbvr"]
