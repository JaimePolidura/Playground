package es.jaime.producer;

import lombok.SneakyThrows;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.clients.producer.RecordMetadata;

import java.util.Properties;
import java.util.concurrent.Future;

public final class SimpleProducer {
    @SneakyThrows
    public static void main(String[] args) {
        Properties props = new Properties();
        props.put("bootstrap.servers", "localhost:9092");
        props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer");
        props.put("value.serializer", "org.apache.kafka.common.serialization.StringSerializer");
        //Solo un mensaje podra estar en "vuelo" a la vez. De esta forma garantizamos el orden de mensajes de envio
        //dentro de un productor
        props.put("flight.requests.per.session=1", "1");

        KafkaProducer<String, String> kafkaProducer = new KafkaProducer(props);

        ProducerRecord<String, String> record = new ProducerRecord<>(
                "CustomerCountry",
                 "Precision Products",
                "France"
        );

        //El futuro se resolvera cuando la petition se envie al broker de kafka
        Future<RecordMetadata> futureRequest = kafkaProducer.send(record);
        futureRequest.get();
        System.out.println("Enviado");

        kafkaProducer.send(record, (recordMetadata, e) -> System.out.println("El mensaje se ha enviado"));
        Thread.sleep(5000);
    }
}
