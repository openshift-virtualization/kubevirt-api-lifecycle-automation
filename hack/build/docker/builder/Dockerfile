FROM quay.io/centos/centos:stream9

ENV GIMME_GO_VERSION=1.22.6
ENV OPERATOR_COURIER_VERSION=2.1.11
ENV GOLANGCI_LINT_VERSION=v1.60.1
ENV GOPATH="/go"
ENV GOBIN="/usr/bin"
ENV GO111MODULE="on"

# Install packages
RUN 	dnf -y install dnf-plugins-core && \
	dnf config-manager --set-enable crb && dnf update -y && \
	dnf install -y \
	cpio \
	diffutils \
	git \
	python3-pip \
	python3-devel \
	gcc \
	gcc-c++ \
	glibc-devel \
	findutils \
	autoconf \
	automake \
	libtool \
	jq \
	rsync-daemon \
	rsync \
	patch \
	libnbd-devel \
	nbdkit-devel \
	unzip \
    protobuf-compiler \
	&& dnf clean all

# reference to master is for an external repo and can't yet be changed
RUN mkdir -p /gimme && curl -sL \
    https://raw.githubusercontent.com/travis-ci/gimme/master/gimme | \
    HOME=/gimme bash >> /etc/profile.d/gimme.sh

# Install persistent go packages
RUN \
	source /etc/profile.d/gimme.sh && \
	eval $(go env) && \
	go install golang.org/x/tools/cmd/goimports@latest && \
	go install mvdan.cc/sh/cmd/shfmt@latest && \
	go install github.com/mattn/goveralls@latest && \
	go install golang.org/x/lint/golint@latest && \
    go install -v github.com/golang/protobuf/protoc-gen-go@1643683 && \
	go install github.com/rmohr/go-swagger-utils/swagger-doc@latest && \
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.14.0 && \
    rm -rf "${GOPATH}/pkg"

RUN set -x && \
    source /etc/profile.d/gimme.sh && \
    go install -v mvdan.cc/gofumpt@v0.6.0 && \
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go env GOROOT)"/bin $GOLANGCI_LINT_VERSION


ADD entrypoint.sh /entrypoint.sh

ADD rsyncd.conf /etc/rsyncd.conf
# Add /root/go/src/github.com/kubevirt/kubevirt-api-lifecycle-automation and mark it as a safe directory for git
RUN mkdir -p /root/go/src/github.com/kubevirt/kubevirt-api-lifecycle-automation && \
    git config --global --add safe.directory /root/go/src/github.com/kubevirt/kubevirt-api-lifecycle-automation

ENTRYPOINT [ "/entrypoint.sh" ]
