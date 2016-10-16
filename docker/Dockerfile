FROM alpine:latest
ENV USER strings
ENV HOME /home/$USER
ENV GOPATH /home/$USER/gocode
ENV PATH /usr/local/via/bin:$GOPATH/bin:/sbin:/bin:/usr/bin:/usr/sbin
RUN adduser -D $USER
RUN mkdir -p /usr/local/via
RUN chown $USER /usr/local/via
RUN apk update
RUN apk add bash grep ncurses git
USER $USER
WORKDIR /home/$USER
ENV EDITOR emacsclient
COPY via /usr/bin
RUN mkdir -p $HOME/publish/repo/master
RUN via sync
RUN via install devel
RUN via install emacs