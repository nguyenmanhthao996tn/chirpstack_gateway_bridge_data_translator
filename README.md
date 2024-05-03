# Chirpstack Gateway Bridge Data Translator

The Chirpstack Gateway Bridge Data Translator utility converts [Chirpstack Gateway Protobuf messages](https://www.chirpstack.io/docs/chirpstack-gateway-bridge/payloads/events.html) in the MQTT Server to JSON messages.

*This repository is developed / tested with Chirpstack V4 (v4.5.1)*

# Caution

The JSON messages could also be archived by setting `json=true` in `[integration.mqtt]` section of the **chirpstack.toml** file. Further details can be found in the [ChirpStack Documentation (v4)](https://www.chirpstack.io/docs/chirpstack/configuration.html).

# Why?

In production, Chirpstack recommends using Protobuf (Protocol Buffers binary encoding) over the JSON messages (for debugging). However, not all programming languages are well-supported with the Protobuf decoding feature. Having a dedicated translating/converting Docker container that runs on the internal network could be your fast & secure solution.

# How to use?

Clone this repository & run:

```
docker-compose up -d
```

Update the configuration file which is located at `config/config.toml` and (re)start the Docker container.

---

###### 🍀🤞 glhf 🤞🍀
