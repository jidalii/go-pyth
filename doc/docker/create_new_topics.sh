#!/bin/bash

#  Kafka container name and port
KAFKA_CONTAINER="pyth-kafka11"
KAFKA_PORT="9192"

# Kafka topic name and configuration
TOPIC_NAME="pythBusiness"
PARTITIONS=2
REPLICATION_FACTOR=3

echo "Creating topic '$TOPIC_NAME' on $KAFKA_CONTAINER (port $KAFKA_PORT)..."

docker exec -it "$KAFKA_CONTAINER" \
    kafka-topics.sh --create \
    --bootstrap-server "pyth-kafka11:9092,pyth-kafka22:9092,pyth-kafka33:9092" \
    --topic "$TOPIC_NAME" \
    --partitions "$PARTITIONS" \
    --replication-factor "$REPLICATION_FACTOR"

if [ $? -eq 0 ]; then
    echo "Topic '$TOPIC_NAME' created successfully on $KAFKA_CONTAINER."
else
    echo "Failed to create topic '$TOPIC_NAME' on $KAFKA_CONTAINER."
fi