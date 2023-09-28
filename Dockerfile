FROM golang:1.20.3-alpine3.17

# Modified from https://github.com/man1207/sipexer

LABEL maintainer="Accruals Unit <tup12806@temple.edu>"
LABEL desription="SIPexer image"
LABEL vcs-type="git"
LABEL vcs-url="https://github.com/gummyfrog/sipexer.git"

COPY . /src

WORKDIR /src

RUN go get ./... \
    && go build .

RUN mv sipexer /usr/local/bin/sipexer

ENTRYPOINT ["sipexer"]

CMD ["--help"]
