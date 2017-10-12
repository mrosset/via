#FROM scratch
FROM busybox
ENV USER mrosset
ENV HOME /home/$USER
ENV PREFIX /tmp/via
ENV GOPATH /home/$USER/gocode
ENV PATH $PREFIX/bin:$GOPATH/bin:/sbin:/bin:/usr/bin:/usr/sbin
COPY bin /bin
COPY etc /etc
COPY usr /usr
RUN ln -s $PREFIX/lib /lib64
RUN mkdir -p /usr/bin
RUN ln -s $PREFIX/bin/file /usr/bin/file
RUN ln -s $PREFIX/bin/env /usr/bin/env
RUN ln -s $PREFIX/bin/bash /usr/bin/bash
WORKDIR /home/$USER
USER $USER
EXPOSE 1122
#CMD ["/home/mrosset/gocode/bin/via","daemon"]
CMD ["/bin/bash","--login"]
