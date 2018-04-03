FROM scratch
MAINTAINER Xu Chenhao <xu.chenhao@hotmail.com>

ADD hyperlook /
EXPOSE 2053
ENTRYPOINT [ "/hyperlook" ]
