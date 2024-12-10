FROM rockylinux:8-minimal

# Use golang 1.22.1
ENV GOLANG_VERSION 1.22.1

# Install required libraries to build discovery image.
# hadolint ignore=DL3041
RUN microdnf install -y git wget tar gcc libguestfs-tools isomd5sum && dnf clean all

# Install golang.
RUN wget -q https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz &&\
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm -rf go${GOLANG_VERSION}.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin

WORKDIR /app

# Download discovery remaster files.
RUN mkdir "remaster" && git clone https://github.com/theforeman/foreman-discovery-image.git remaster/.
ENV DISCOVERY_REMASTER remaster/aux/remaster/discovery-remaster

# Download the Discovery fdi bootable image.
RUN mkdir "fdi" && wget -q https://downloads.theforeman.org/discovery/releases/4.1/fdi-4.1.0-24d62de.iso -O fdi/fdi-4.1.0.iso
ENV DISCOVERY_BASE_IMAGE fdi/fdi-4.1.0.iso

# Create the images folder.
ENV IMAGES_PATH images
RUN mkdir ${IMAGES_PATH}

# Download Go modules and Build
COPY . .
RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o /fpi

EXPOSE 8080

CMD ["/fpi"]

