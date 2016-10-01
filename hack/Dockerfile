from docker:1.11.2

RUN apk -U add curl git

COPY sshfront /sshfront

COPY auth.sh /auth.sh

COPY handler.sh /handler.sh

VOLUME /workdir

ENTRYPOINT ["/sshfront", "-a", "/auth.sh", "/handler.sh"]
