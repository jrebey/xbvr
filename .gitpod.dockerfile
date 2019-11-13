FROM gitpod/workspace-full

ENV HOME=/home/gitpod
WORKDIR $HOME
USER gitpod

RUN pip3 install sqlite-web

ENV GO_VERSION=1.12 \
  GOPATH=$HOME/go-packages \
  GOROOT=$HOME/go
RUN export PATH=$(echo "$PATH" | sed -e 's|:/workspace/go/bin||' -e 's|:/home/gitpod/go/bin||' -e 's|:/home/gitpod/go-packages/bin||')
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH

RUN go get -u -v \
  github.com/UnnoTed/fileb0x && \
  # Temp workaround for broken modd deps
  # github.com/cortesi/modd/cmd/modd && \
  git clone https://github.com/cortesi/modd && \
  cd modd && \
  go get mvdan.cc/sh@8aeb0734cd0f && \
  go install ./cmd/modd && \
  sudo rm -rf $GOPATH/src && \
  sudo rm -rf $GOPATH/pkg
# user Go packages
ENV GOPATH=/workspace/go \
  PATH=/workspace/go/bin:$PATH

RUN pip3 install --no-cache-dir cython && \
  pip3 install --no-cache-dir flask peewee sqlite-web

RUN git config --global alias.gofmt \
  '!echo $(git diff --cached --name-only --diff-filter=ACM | grep '.go$') | \
  xargs gofmt -w -l | \
  xargs git add'

USER root
