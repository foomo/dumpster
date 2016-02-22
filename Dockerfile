FROM scratch

RUN mkdir -p /bin
ADD /go/bin/dumpster /bin/dumpster

ENTRYPOINT ["/bin/dumpster"]
