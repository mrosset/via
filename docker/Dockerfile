#FROM scratch
FROM busybox
ENV USER strings
ENV HOME /home/$USER
ENV GOPATH /home/$USER/gocode
ENV PATH /usr/local/via/bin:$GOPATH/bin:/sbin:/bin:/usr/bin:/usr/sbin
COPY etc /etc
COPY usr /usr
RUN mkdir -p /usr/local/via
RUN chown -R strings /usr/local/via
RUN ln -s /usr/local/via/lib /lib64
WORKDIR /home/$USER
USER $USER
EXPOSE 1122
CMD ["/usr/local/via/bin/bash","--login"]
#CMD ["/usr/bin/circuit","start","-a","172.17.0.2:1122"]
