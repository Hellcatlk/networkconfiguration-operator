FROM  alpine:latest

# Copy file
WORKDIR /
COPY ./bin/manager .
COPY ./bin/network-runner /usr/bin

# Prepare running environment
RUN apk add ansible openssh sshpass py3-pip gcc g++ --no-cache && \
    apk add python3-dev libc-dev linux-headers --no-cache && \
    pip3 install networking-ansible && \
    cp -rf /usr/lib/python3.9/site-packages/etc/ansible /etc/ && \
    ansible-galaxy collection install fujitsu.fos && \
    ln -s ~/.ansible/collections/ansible_collections/fujitsu/fos/plugins ~/.ansible/ && \
    echo "StrictHostKeyChecking no" >> /etc/ssh/ssh_config

ENTRYPOINT ["/manager"]
