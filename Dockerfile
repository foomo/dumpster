FROM scratch

ADD bin/dumpster /bin/dumpster

ENTRYPOINT ["/bin/dumpster"]
