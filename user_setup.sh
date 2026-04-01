#!/bin/bash
for user in matt tanner aidan
do
    useradd -m $user
    echo -e "ubuntu\nubuntu" | passwd $user
    sudo passwd -e $user
    usermod -aG sudo $user
    usermod -aG devs $user
    mkdir -p /home/$user/
    ln -s /mnt/capstone /home/$user/
done