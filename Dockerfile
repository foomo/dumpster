FROM scratch

ADD $GOPATH/bin/dumpster /bin/dumpster

ENTRYPOINT ["/bin/dumpster"]
