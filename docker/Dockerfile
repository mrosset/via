#FROM scratch
FROM busybox
ENV USER mrosset
ENV HOME /home/$USER
ENV GOPATH /home/$USER/gocode
ENV PATH $HOME/via/bin:$GOPATH/bin:/sbin:/bin:/usr/bin:/usr/sbin
COPY bin /bin
COPY etc /etc
COPY usr /usr
RUN ln -s $HOME/via/lib /lib64
RUN mkdir -p /usr/bin
RUN mkdir -p $HOME/via
RUN chown $USER.$USER $HOME/via
RUN ln -s $HOME/via/bin/file /usr/bin/file
WORKDIR /home/$USER
USER $USER
EXPOSE 1122
CMD ["/bin/bash","--login"]
#CMD ["/usr/bin/circuit","start","-a","172.17.0.2:1122"]
