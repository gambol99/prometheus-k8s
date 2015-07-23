#
#   Author: Rohith
#   Date: 2015-07-23 12:00:34 +0100 (Thu, 23 Jul 2015)
#
#  vim:ts=2:sw=2:et
#
FROM gliderlabs/alpine:3.1
MAINTAINER Rohith <gambol99@gmail.com>

ADD bin/prometheus-k8s /bin/prometheus-k8s
RUN chmod +x /bin/prometheus-k8s

ENTRYPOINT [ "/bin/prometheus-k8s" ]
