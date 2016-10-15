FROM alpine:latest
ENV USER via
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
COPY via /usr/bin