MQTT Practice Materials
=======================

See also [related slides](https://lectures-dot-kpi-architecture-course.appspot.com/architecture/15-mqtt-practice.slide).

After cloning the repo, run with
```bash
docker-compose up
```

You should be able to connect to the MQTT broker using `localhost` address.
For instance,
```bash
mosquitto_sub -h localhost -t /test/inception
```
