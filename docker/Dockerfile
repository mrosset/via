FROM alpine:latest
ENV USER strings
ENV HOME /home/$USER
ENV GOPATH /home/$USER/gocode
ENV PATH /usr/local/via/bin:$GOPATH/bin:/sbin:/bin:/usr/bin:/usr/sbin
RUN apk update
RUN apk add sudo grep ncurses bash ca-certificates
RUN adduser -D $USER
RUN echo "$USER ALL=(ALL) NOPASSWD: ALL" > "/etc/sudoers.d/$USER"
#COPY usr /usr
COPY usr/bin/via /usr/bin/via
RUN mkdir -p /usr/local/via
RUN chown -R $USER /usr/local/via
WORKDIR /home/$USER
USER $USER