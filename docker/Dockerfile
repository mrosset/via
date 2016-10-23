FROM alpine:latest
ENV USER strings
ENV HOME /home/$USER
ENV GOPATH /home/$USER/gocode
ENV PATH /usr/local/via/bin:$GOPATH/bin:/sbin:/bin:/usr/bin:/usr/sbin
RUN apk update
RUN apk add sudo grep ncurses git bash
RUN adduser -D $USER
RUN mkdir -p /usr/local/via
RUN chown $USER /usr/local/via
RUN echo "$USER ALL=(ALL) NOPASSWD: ALL" > "/etc/sudoers.d/$USER"
COPY via /usr/bin
WORKDIR /home/$USER
USER $USER