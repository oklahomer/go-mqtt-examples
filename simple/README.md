This example shows how to publish and subscribe to a given topic.

# MQTT Broker activation
```
$ docker-compose up -d
Starting simple_mqtt_1 ... done
```

Log file is located at `mosquitto/log/mosquitto.log`.
See `mosquitto/config/mosquitto.conf` for all configuration.

To finish its operation, run below command to stop.
```
$ docker-compose stop
Stopping simple_mqtt_1 ... done
```

# Subscription
After starting MQTT server, execute below command to start subscription.
As long as the publication script is running, received message follows.
```
$ go ./sub/main.go
2020/10/18 14:43:35 Connected
```

# Publication
After starting MQTT server, execute below command to start publication.
A message is published every 3 seconds.
```
$ go run ./pub/main.go 
2020/10/18 14:44:13 Connected
2020/10/18 14:44:16 Published
```