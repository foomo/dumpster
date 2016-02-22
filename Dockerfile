FROM scratch

ADD /go/bin/dumpster /bin/dumpster

ENTRYPOINT ["/bin/dumpster"]
