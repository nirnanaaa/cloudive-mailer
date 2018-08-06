# Distributed Mailing

```sh
# start up kafka

docker run --rm -ti -p 9009:9009 -e CLOUDIVE_KAFKA_BROKERS=kafka:9092 cloudive/mailer master
docker run --rm -ti -e CLOUDIVE_KAFKA_BROKERS=kafka:9092 cloudive/mailer slave
```