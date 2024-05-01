package main

import (
    "fmt"
    mqtt "github.com/eclipse/paho.mqtt.golang"
    "log"
    "time"
    "github.com/golang/protobuf/proto"
    "strings"
    "github.com/chirpstack/chirpstack/api/go/v4/gw"
    "encoding/json"
	"os"
	"github.com/BurntSushi/toml"
)

// ==================================

type Config struct {
	Mqtt_broker_address 		string
	Mqtt_broker_port 				int
	Mqtt_client_id 					string
	Mqtt_client_username 		string
	Mqtt_client_password 		string
	Mqtt_listen_topic 			string
	Mqtt_translate_topic 		string
}
var conf Config

// ==================================

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

    split_topic_arr := strings.Split(string(msg.Topic()), "/")
    if ((split_topic_arr[3] == "event") && (split_topic_arr[4] == "up")) {
        fmt.Println("====================");
        fmt.Println("Topic: event/up");

        gw_uplink := &gw.UplinkFrame{}
        err := proto.Unmarshal(msg.Payload(), gw_uplink)
        if err != nil {
            log.Fatal("unmarshaling error: ", err)
        }

        gw_data_json, err := json.Marshal(gw_uplink)
        fmt.Println("JSON: ", string(gw_data_json));

        send_translated_message(client, msg.Topic(), string(gw_data_json));

        fmt.Println("====================");

    } else if ((split_topic_arr[3] == "event") && (split_topic_arr[4] == "stats")) {
        fmt.Println("====================");
        fmt.Println("Topic: event/stats");

        gw_data := &gw.GatewayStats{}
        err := proto.Unmarshal(msg.Payload(), gw_data)
        if err != nil {
            log.Fatal("unmarshaling error: ", err)
        }

        gw_data_json, err := json.Marshal(gw_data)
        fmt.Println("JSON: ", string(gw_data_json));

        send_translated_message(client, msg.Topic(), string(gw_data_json));

        fmt.Println("====================");

    } else if ((split_topic_arr[3] == "event") && (split_topic_arr[4] == "ack")) {
        fmt.Println("====================");
        fmt.Println("Topic: event/ack");

        gw_data := &gw.DownlinkTxAck{}
        err := proto.Unmarshal(msg.Payload(), gw_data)
        if err != nil {
            log.Fatal("unmarshaling error: ", err)
        }

        gw_data_json, err := json.Marshal(gw_data)
        fmt.Println("JSON: ", string(gw_data_json));

        send_translated_message(client, msg.Topic(), string(gw_data_json));

        fmt.Println("====================");

    } else if ((split_topic_arr[3] == "event") && (split_topic_arr[4] == "exec")) {
        fmt.Println("====================");
        fmt.Println("Topic: event/exec");

        gw_data := &gw.GatewayCommandExecResponse{}
        err := proto.Unmarshal(msg.Payload(), gw_data)
        if err != nil {
            log.Fatal("unmarshaling error: ", err)
        }

        gw_data_json, err := json.Marshal(gw_data)
        fmt.Println("JSON: ", string(gw_data_json));

        send_translated_message(client, msg.Topic(), string(gw_data_json));

        fmt.Println("====================");

        } else if ((split_topic_arr[3] == "event") && (split_topic_arr[4] == "raw")) {
            fmt.Println("====================");
            fmt.Println("Topic: event/raw");
    
            gw_data := &gw.RawPacketForwarderEvent{}
            err := proto.Unmarshal(msg.Payload(), gw_data)
            if err != nil {
                log.Fatal("unmarshaling error: ", err)
            }
    
            gw_data_json, err := json.Marshal(gw_data)
            fmt.Println("JSON: ", string(gw_data_json));

            send_translated_message(client, msg.Topic(), string(gw_data_json));
    
            fmt.Println("====================");

            } else if ((split_topic_arr[3] == "command") && (split_topic_arr[4] == "down")) {
                fmt.Println("====================");
                fmt.Println("Topic: command/down");
        
                gw_data := &gw.DownlinkFrame{}
                err := proto.Unmarshal(msg.Payload(), gw_data)
                if err != nil {
                    log.Fatal("unmarshaling error: ", err)
                }
        
                gw_data_json, err := json.Marshal(gw_data)
                fmt.Println("JSON: ", string(gw_data_json));
    
                send_translated_message(client, msg.Topic(), string(gw_data_json));
        
                fmt.Println("====================");
        
    } else {
        fmt.Println("Message type/topic handler is not defined yet, topic: ", msg.Topic());
    }
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
    fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
    fmt.Println("Connect lost: ", err)
}

func main() {
    set_conf_default_value(&conf)
    read_conf_from_toml(&conf)

    var broker = conf.Mqtt_broker_address
    var port = conf.Mqtt_broker_port
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
    opts.SetClientID(conf.Mqtt_client_id)
    opts.SetUsername(conf.Mqtt_client_username)
    opts.SetPassword(conf.Mqtt_client_password)
    opts.SetDefaultPublishHandler(messagePubHandler)
    opts.OnConnect = connectHandler
    opts.OnConnectionLost = connectLostHandler
    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    sub(client)
    // publish(client)

    for {}

    client.Disconnect(250)
}

func publish(client mqtt.Client) {
    num := 10
    for i := 0; i < num; i++ {
        text := fmt.Sprintf("Message %d", i)
        token := client.Publish("topic/test", 0, false, text)
        token.Wait()
        time.Sleep(time.Second)
    }
}

func send_translated_message(client mqtt.Client, topic string, text string) {
    translated_topic := strings.Join([]string{conf.Mqtt_translate_topic, topic}, "/")
    translated_topic = strings.ReplaceAll(translated_topic, "//", "/")
    fmt.Println("Sending translated JSON data to ", translated_topic);

    token := client.Publish(translated_topic, 0, false, text)
    token.Wait()
}

func sub(client mqtt.Client) {
    listen_topic := conf.Mqtt_listen_topic
    listen_topic = strings.ReplaceAll(listen_topic, "//", "/")

    topic := listen_topic
    token := client.Subscribe(topic, 1, nil)
    token.Wait()
    fmt.Println("Subscribed to topic: %s", topic)
}

func set_conf_default_value(conf *Config) {
	conf.Mqtt_broker_address = "manhthao.ovh"
	conf.Mqtt_broker_port = 11883
	conf.Mqtt_client_id = "client_id"
	conf.Mqtt_client_username = ""
	conf.Mqtt_client_password = ""
	conf.Mqtt_listen_topic = "as923_2/gateway/7276ff000f061f16/"
	conf.Mqtt_translate_topic = "translated/"
}

func read_conf_from_toml(conf *Config) {
    dat, err := os.ReadFile("./config.toml")
	check(err)

    _, err1 := toml.Decode(string(dat), &conf)
	check(err1)

    fmt.Println("Mqtt_broker_address: ", conf.Mqtt_broker_address)
	fmt.Println("Mqtt_broker_port: ", conf.Mqtt_broker_port)
	fmt.Println("Mqtt_client_id: ", conf.Mqtt_client_id)
	fmt.Println("Mqtt_client_username: ", conf.Mqtt_client_username)
	fmt.Println("Mqtt_client_password: ", conf.Mqtt_client_password)
	fmt.Println("Mqtt_listen_topic: ", conf.Mqtt_listen_topic)
	fmt.Println("Mqtt_translate_topic: ", conf.Mqtt_translate_topic)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}