FROM scratch
ENV USER strings
ENV HOME /home/$USER
ENV GOPATH /home/$USER/gocode
ENV PATH /usr/local/via/bin:$GOPATH/bin:/sbin:/bin:/usr/bin:/usr/sbin
COPY bin /bin
COPY usr /usr
COPY etc /etc
#RUN mkdir -p /usr/local/via
#RUN chown -R $USER /usr/local/via
WORKDIR /home/$USER
USER $USER
CMD ["/usr/bin/via","start"]