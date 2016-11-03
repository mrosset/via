#FROM scratch
FROM busybox
ENV USER strings
ENV HOME /home/$USER
ENV GOPATH /home/$USER/gocode
ENV PATH /usr/local/via/bin:$GOPATH/bin:/sbin:/bin:/usr/bin:/usr/sbin
COPY home /home
COPY etc /etc
COPY usr /usr
WORKDIR /home/$USER
USER $USER
CMD ["/usr/bin/via","start"]