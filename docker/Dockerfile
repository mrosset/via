#FROM scratch
FROM busybox
ENV USER mrosset
ENV HOME /home/$USER
ENV PREFIX $HOME/via
ENV GOPATH /home/$USER/gocode
ENV PATH $PREFIX/bin:$HOME/via/bin:$GOPATH/bin:/sbin:/bin:/usr/bin:/usr/sbin
COPY bin /bin
COPY etc /etc
RUN mkdir -p $PREFIX
RUN chown -R $USER $HOME
RUN ln -s $PREFIX/lib /lib64
RUN ln -s $PREFIX/include /usr/
RUN mkdir -p /usr/bin
RUN ln -s $PREFIX/bin/file /usr/bin/file
RUN ln -s $PREFIX/bin/env /usr/bin/env
RUN ln -s $PREFIX/bin/bash /usr/bin/bash
WORKDIR /home/$USER
USER $USER
EXPOSE 1122
CMD ["/usr/bin/via","daemon"]
#CMD ["/bin/bash","--login"]
COPY usr/bin/via /usr/bin/