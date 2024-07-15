package es.jaime.producer.customserializer;

import lombok.SneakyThrows;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerRecord;

import java.util.Properties;

public final class ProducerWithCustomSerializer {
    @SneakyThrows
    public static void main(String[] args) {
        ProducerRecord<String, Customer> record = new ProducerRecord<>(
                "Customers",
                "1",
                Customer.builder()
                        .customerId(1)
                        .name("Jaime")
                        .build()
        );
        Properties props = new Properties();
        props.put("bootstrap.servers", "localhost:9092");
        props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer");
        props.put("value.serializer", CustomerSerializer.class.getName());

        KafkaProducer<String, Customer> kafkaProducer = new KafkaProducer<>(props);
        kafkaProducer.send(record).get();
        System.out.println("El otro mensaje se ha enviado");
    }
}
