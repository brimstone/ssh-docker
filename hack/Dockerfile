from docker:1.11.2

RUN apk -U add curl git openssh-client

RUN apk -U add -t .build go \
 && export GOPATH=/tmp/go \
 && go get -v github.com/gliderlabs/sshfront \
 && mv /tmp/go/bin/* /usr/bin/ \
 && rm -rf /tmp/go \
 && apk del .build

RUN git config --global user.email "brimbot@the.narro.ws" \
 && git config --global user.name "Brimbot https://github.com/brimstone/ssh-docker"

COPY auth.sh /auth.sh

COPY handler.sh /handler.sh

COPY loader /loader

ENTRYPOINT ["/loader"]
