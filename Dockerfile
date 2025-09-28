FROM ubuntu:latest
LABEL authors="dryundel"

ENTRYPOINT ["top", "-b"]