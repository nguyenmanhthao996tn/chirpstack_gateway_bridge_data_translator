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
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
    // fmt.Printf("Received message: %v from topic: %s\n", msg.Payload(), msg.Topic())

    split_topic_arr := strings.Split(string(msg.Topic()), "/")
    // fmt.Printf("split_topic_arr: ", split_topic_arr)
    if ((split_topic_arr[3] == "event") && (split_topic_arr[4] == "up")) {
        fmt.Println("====================");
        fmt.Printf("Topic: event/up\n");

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
        fmt.Printf("Topic: event/stats\n");

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
        fmt.Printf("Topic: event/ack\n");

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
        fmt.Printf("Topic: event/exec\n");

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
            fmt.Printf("Topic: event/raw\n");
    
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
                fmt.Printf("Topic: command/down\n");
        
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
        fmt.Printf("Message type / topic handler is not defined yet, topic: %s\n", msg.Topic());
    }
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
    fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
    fmt.Printf("Connect lost: %v", err)
}

func main() {
    var broker = "manhthao.ovh"
    var port = 11883
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
    opts.SetClientID("go_mqtt_client")
    opts.SetUsername("emqx")
    opts.SetPassword("public")
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
    translated_topic := strings.Join([]string{"translated", topic}, "/")
    fmt.Println("Sending translated JSON data to ", translated_topic);

    token := client.Publish(translated_topic, 0, false, text)
    token.Wait()
}

func sub(client mqtt.Client) {
    topic := "as923_2/gateway/7276ff000f061f16/#"
    token := client.Subscribe(topic, 1, nil)
    token.Wait()
    fmt.Println("Subscribed to topic: %s", topic)
}