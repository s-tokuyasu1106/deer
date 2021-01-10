FROM centos:7

# set timezone
RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# compile and install vim
RUN rpm --rebuilddb && yum update -y && \
    yum -y install git gcc ncurses-devel make && \
    yum -y install httpd which && \
    yum -y install vim && \
    rm -rf /var/cache/yum/* && \
    yum clean all

FROM golang:latest

# コンテナ作業ディレクトリの変更
WORKDIR /go/src/deer/go/chat/

ENV GO111MODULE=auto

# パッケージのインストール
RUN go get github.com/lib/pq

# ホストOSの ./go の中身を作業ディレクトリに追加
ADD ./go .

# Dockerfileファイル中に以下のような記述を追加し、8080番ポートを公開する
EXPOSE 8080
RUN go get golang.org/x/tools/cmd/...
# RUN go get github.com/pilu/fresh
RUN go get -u golang.org/x/lint/golint
RUN go get github.com/oxequa/realize
RUN go get github.com/gorilla/websocket
#CMD ["fresh"]
# CMD "fresh"
# Air インストール
RUN go get -u github.com/cosmtrek/air
COPY go.mod go.sum ./
RUN go mod download
# コンテナ実行時のデフォルトを設定する
# ライブリロードを実行する





# ENV GO_INSTALL_DIR /go/src/deer/go
# ENV HOME /root
# ENV GO_VERSION 1.12.4

# # set timezone
# RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# # compile and install vim
# RUN rpm --rebuilddb && yum update -y && \
#     yum -y install git gcc ncurses-devel make && \
#     yum -y install httpd which && \
#     rm -rf /var/cache/yum/* && \
#     yum clean all

# WORKDIR /go/src/deer/go
# RUN git clone https://github.com/vim/vim.git && \
#     cd vim && \
#     make && make install && make clean && \
#     rm -rf vim

# # Your .vimrc
# # ADD .vimrc ~/

# # get golang and set $GOPATH
# WORKDIR ${GO_INSTALL_DIR}

# RUN curl https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz > go.tar.gz && \
#     tar zxf go.tar.gz && \
#     rm go.tar.gz

# # RUN go get github.com/lib/pq
# # RUN go get github.com/oxequa/realize

# ADD ./go .

# WORKDIR ${HOME}
# RUN echo 'export GOROOT=/usr/local/go' >> ${HOME}/.bashrc
# RUN echo 'export GOPATH=${HOME}/go' >> ${HOME}/.bashrc
# RUN echo 'export PATH=${GOROOT}/bin:${PATH}' >> ${HOME}/.bashrc
# RUN mkdir -p ${HOME}/go/src

# # make directory to share with host
# RUN mkdir -p /share

# CMD [ "/sbin/init" ]