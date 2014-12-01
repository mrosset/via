FROM strings/via
ENV PATH /usr/local/via/bin
ENV EDITOR vim
ADD via/via /usr/local/via/bin/
USER strings
