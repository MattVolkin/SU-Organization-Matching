#!/bin/bash
set -euo pipefail

SHARED_ROOT="/mnt/capstone"
REPO_PATH="/mnt/capstone/SU-Organization-Matching"
SHARED_GROUP="devs"

if ! getent group "$SHARED_GROUP" >/dev/null; then
    groupadd "$SHARED_GROUP"
fi

for user in matt tanner aidan ben
do
    if ! id -u "$user" >/dev/null 2>&1; then
        useradd -m "$user"
        echo -e "ubuntu\nubuntu" | passwd "$user"
        passwd -e "$user"
    fi

    usermod -aG sudo "$user"
    usermod -aG "$SHARED_GROUP" "$user"
    mkdir -p "/home/$user"
    ln -sfn "$SHARED_ROOT" "/home/$user/capstone"

    # Make new files created by these users group-writable by default.
    grep -qxF 'umask 0002' "/home/$user/.profile" || echo 'umask 0002' >> "/home/$user/.profile"
done

# Enforce shared ownership/permissions and inheritance on the shared tree.
chgrp -R "$SHARED_GROUP" "$SHARED_ROOT"
find "$SHARED_ROOT" -type d -exec chmod 2775 {} +
find "$SHARED_ROOT" -type f -exec chmod g+rw {} +

setfacl -R -m g:"$SHARED_GROUP":rwx "$SHARED_ROOT"
setfacl -R -d -m g:"$SHARED_GROUP":rwx "$SHARED_ROOT"

# Git repositories work best without executable ACLs on regular files.
find "$REPO_PATH" -type f -exec chmod g-x {} +