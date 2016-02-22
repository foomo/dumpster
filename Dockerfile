FROM scratch

ADD dumpster /bin/dumpster

ENTRYPOINT ["/bin/dumpster"]
