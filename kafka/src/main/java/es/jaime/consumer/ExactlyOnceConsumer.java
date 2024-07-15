package es.jaime.consumer;

import org.apache.kafka.clients.consumer.ConsumerRebalanceListener;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.clients.consumer.KafkaConsumer;
import org.apache.kafka.common.TopicPartition;

import java.time.Duration;
import java.time.temporal.ChronoUnit;
import java.util.Collection;
import java.util.List;
import java.util.Properties;

public final class ExactlyOnceConsumer {
    public static final KafkaConsumer<String, String> consumer = new KafkaConsumer<>(getProperties());

    public static void main(String[] args) {
        consumer.subscribe(List.of("customerCountries"));
        //Nos aseguramos de que nos conectamos a una particion
        consumer.poll(0);

        //Por cada particion, marcamos que tenemos que empezar a consumir desde cierto offset
        for (TopicPartition partition: consumer.assignment()) {
            consumer.seek(partition, getOffsetFromDB(partition));
        }

        while (true) {
            for (var record : consumer.poll(Duration.of(100, ChronoUnit.MILLIS))) {
                processRecord(record);
                storeInDb(record);
            }

            commitDBTransaction();
        }
    }

    private static void storeInDb(ConsumerRecord<String, String> record) {
        //A parte de guardar el record en la BDD, guardamos también el offset
    }

    private static void processRecord(ConsumerRecord<String, String> record) {}

    public static class SaveOffsetsOnRebalance implements ConsumerRebalanceListener {
        public void onPartitionsRevoked(Collection<TopicPartition> partitions) {
            commitDBTransaction();
        }

        public void onPartitionsAssigned(Collection<TopicPartition>
                                                 partitions) {
            for(TopicPartition partition: partitions) {
                consumer.seek(partition, getOffsetFromDB(partition));
            }
        }
    }

    private static long getOffsetFromDB(TopicPartition partition) {
        return 0;
    }

    private static Properties getProperties() {
        Properties props = new Properties();
        props.put("bootstrap.servers", "localhost:9092");
        props.put("key.serializer", "org.apache.kafka.common.serialization.StringSerializer");
        props.put("value.serializer", "org.apache.kafka.common.serialization.StringSerializer");
        props.put("enable.auto.commit", "false");
        return props;
    }

    private static void commitDBTransaction() {

    }
}
