FROM alpine:latest
ENV USER strings
ENV HOME /home/$USER
ENV GOPATH /home/$USER/gocode
ENV PATH /usr/local/via/bin:$GOPATH/bin:/sbin:/bin:/usr/bin:/usr/sbin
RUN adduser -D $USER
RUN mkdir -p /usr/local/via
RUN chown $USER /usr/local/via
RUN apk update
RUN apk add bash grep ncurses git sudo
RUN echo "$USER ALL=(ALL) NOPASSWD: ALL" > "/etc/sudoers.d/$USER"
USER $USER
WORKDIR /home/$USER
#COPY via /usr/bin
ENV EDITOR emacsclient
#RUN via sync
#RUN via install devel
