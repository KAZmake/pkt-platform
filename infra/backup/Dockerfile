FROM alpine:3.20

# pg_dump + pg_restore + psql
RUN apk add --no-cache postgresql16-client gzip bash curl

# rclone
RUN curl -fsSL https://rclone.org/install.sh | bash

COPY backup.sh  /usr/local/bin/backup
COPY restore.sh /usr/local/bin/restore
RUN chmod +x /usr/local/bin/backup /usr/local/bin/restore

# rclone config is injected at runtime via env vars or mounted secret
CMD ["backup"]
