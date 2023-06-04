#Deriving the latest base image
FROM docker.io/python:3.12.0b1-alpine3.18

#Labels as key value pair
LABEL Maintainer="nvalembois@live.com"
LABEL org.opencontainers.image.source="https://github.com/nvalembois/gh-pr-auto-approver"

# Any working directory can be chosen as per choice like '/' or '/home' etc
# i have chosen /usr/app/src
WORKDIR /usr/app/src

#to COPY the remote file at working directory in container
COPY requirements.txt pullrequest_auto_approver.py ./

# Install pip required library
RUN set -eux; \
	apk add --no-cache --virtual .build-deps \
		gnupg \
		tar \
		xz \
		bluez-dev \
		bzip2-dev \
		dpkg-dev dpkg \
		expat-dev \
		findutils \
		gcc \
		gdbm-dev \
		libc-dev \
		libffi-dev \
		libnsl-dev \
		libtirpc-dev \
		linux-headers \
		make \
		ncurses-dev \
		openssl-dev \
		pax-utils \
		readline-dev \
		sqlite-dev \
		tcl-dev \
		tk \
		tk-dev \
		util-linux-dev \
		xz-dev \
		zlib-dev; \
    PYTHONDONTWRITEBYTECODE=1 \
        pip3 install \
            --no-cache-dir --no-compile \
            -r requirements.txt; \
	apk del --no-network .build-deps

#CMD instruction should be used to run the software
#contained by your image, along with any arguments.

CMD [ "python", "./pullrequest_auto_approver.py"]
