#!/bin/ash

set -e

# Extract the required variables from the environment
ADMIN_TOKEN=$(cat $DOCKER_INFLUXDB_INIT_ADMIN_TOKEN_FILE)
REMOTE_ORG_ID=$(influx org ls --host $REMOTE_URL -t $ADMIN_TOKEN -n $DOCKER_INFLUXDB_INIT_ORG  | awk -v pattern=$DOCKER_INFLUXDB_INIT_ORG '$0 ~ pattern {print $1}')
LOCAL_BUCKET_ID=$(influx bucket ls -t $ADMIN_TOKEN | awk -v pattern=$DOCKER_INFLUXDB_INIT_BUCKET '$0 ~ pattern {print $1}')

# Create a remote on the "source" machine, pointing to the replica.
influx remote create \
  --name replica \
  -o $DOCKER_INFLUXDB_INIT_ORG \
  -t $ADMIN_TOKEN \
  --remote-url $REMOTE_URL \
  --remote-api-token $ADMIN_TOKEN \
  --remote-org-id $REMOTE_ORG_ID --allow-insecure-tls

REMOTE_ID=$(influx remote list -t $ADMIN_TOKEN | awk -v pattern=$REMOTE_URL '$0 ~ pattern {print $1}')

# Create the replication to the source machine.
influx replication create \
  --name $REPLNAME \
  --remote-id $REMOTE_ID \
  --local-bucket-id $LOCAL_BUCKET_ID \
  --remote-bucket $DEST_BUCKET