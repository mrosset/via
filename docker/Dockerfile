#FROM scratch
FROM busybox
ENV USER mrosset
ENV HOME /home/$USER
ENV GOPATH /home/$USER/gocode
ENV PATH /via/bin:$GOPATH/bin:/sbin:/bin:/usr/bin:/usr/sbin
COPY bin /bin
COPY etc /etc
COPY usr /usr
RUN ln -s /via/lib /lib64
RUN mkdir -p /usr/bin
RUN mkdir -p /via
RUN chown $USER.$USER /via
RUN ln -s /via/bin/file /usr/bin/file
WORKDIR /home/$USER
USER $USER
EXPOSE 1122
#CMD ["/home/mrosset/gocode/bin/via","daemon"]
CMD ["/bin/bash","--login"]
