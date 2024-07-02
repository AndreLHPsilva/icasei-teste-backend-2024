#!/bin/sh

echo "Aguardando Kafka iniciar..."
while ! nc -z kafka 29092; do
  sleep 1
done

# Cria o tópico
kafka-topics --bootstrap-server kafka:29092 --create --topic go-to-rails --partitions 1 --replication-factor 1 || true

echo "Tópico criado com sucesso ou já existente"
