#FROM scratch
FROM busybox
ENV USER strings
ENV HOME /home/$USER
ENV GOPATH /home/$USER/gocode
ENV PATH /usr/local/via/bin:$GOPATH/bin:/sbin:/bin:/usr/bin:/usr/sbin
COPY etc /etc
COPY usr /usr
RUN chown -R strings /usr/local/via
WORKDIR /home/$USER
USER $USER
CMD ["/usr/local/via/bin/bash","--login","-o","vi"]
