#!/bin/bash
for user in matt tanner aidan
do
    useradd -m $user
    echo -e "ubuntu\nubuntu" | passwd $user
    mkdir -p /home/$user/capstone
    ln -s /mnt/capstone /home/$user/
done