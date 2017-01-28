FROM centurylink/ca-certs
EXPOSE 3001
COPY flyfox /
ENTRYPOINT ["/flyfox"]
