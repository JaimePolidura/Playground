package es.jaime.consumer;

import org.apache.kafka.clients.consumer.*;
import org.apache.kafka.common.TopicPartition;

import java.time.Duration;
import java.time.temporal.ChronoUnit;
import java.util.*;

public final class SimpleConsumer {
    public static void main(String[] args) {
        Properties props = new Properties();
        props.put("bootstrap.servers", "localhost:9092");
        props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer");
        props.put("value.serializer", "org.apache.kafka.common.serialization.StringSerializer");
        props.put("group.id", "CountryCounter");
        //Cada 5 s el consumidor enviará a kafka el último offset que ha consumido. Si el consumidor crasea, el nuevo consumidor
        //probablemente reprocesara de nuevo los nuevos mensajes.
        props.put("enable.auto.commit", "true");
        //Si es false, a la hora de consumer mensajes, podremos seleccionar cuando vamos a comitear un offset, con commitSync()
        props.put("enable.auto.commit", "false");

        KafkaConsumer<String, String> consumer = new KafkaConsumer<>(props);

        consumer.subscribe(List.of("customerCountries"));

        while (true) {
            for (var record : consumer.poll(Duration.of(100, ChronoUnit.MILLIS))) {
                System.out.printf("Received record of key %s and value %s", record.key(), record.value());
            }
            //Solo cuando "enable.auto.commit" es "false". Si falla se reintenta
            consumer.commitSync();
            //El commit se hace en otro hilo. Si falla, no se reintenta
            //consumer.commitAsync();
        }
    }

    private class HandleRebalance implements ConsumerRebalanceListener {
        public void onPartitionsAssigned(Collection<TopicPartition> partitions) {
            //Cada vez que se asignan nuevas particiones
        }
        public void onPartitionsRevoked(Collection<TopicPartition> partitions) {
            //Cada vez que se asignan particiones
        }
    }
}
