package es.jaime.streams;

import lombok.SneakyThrows;
import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.KafkaStreams;
import org.apache.kafka.streams.KeyValue;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.StreamsConfig;
import org.apache.kafka.streams.kstream.KStream;
import org.apache.kafka.streams.kstream.Named;

import java.util.Arrays;
import java.util.Properties;
import java.util.regex.Pattern;

public final class WordCount {
    @SneakyThrows
    public static void main(String[] args) {
        Properties props = new Properties();
        props.put(StreamsConfig.APPLICATION_ID_CONFIG, "wordcount");
        props.put(StreamsConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9092");
        props.put(StreamsConfig.DEFAULT_KEY_SERDE_CLASS_CONFIG, Serdes.String().getClass().getName());
        props.put(StreamsConfig.DEFAULT_VALUE_SERDE_CLASS_CONFIG, Serdes.String().getClass().getName());

        StreamsBuilder builder = new StreamsBuilder();
        KStream<String, String> sourceTopic = builder.stream("wordcount-input-topic");
        Pattern pattern = Pattern.compile("\\W+");

        KStream<Object, String> counts = sourceTopic
                .flatMapValues(value-> Arrays.asList(pattern.split(value.toLowerCase())))
                .map((key, value) -> new KeyValue<Object, Object>(value, value))
                .filter((key, value) -> !value.equals("the"))
                .groupByKey()
                .count(Named.as("CountStore")).mapValues(value-> Long.toString(value)).toStream();
        counts.to("wordcount-output");

        //Run for 10 seconds
        KafkaStreams streams = new KafkaStreams(builder.build(), props);
        streams.start();
        Thread.sleep(10000);
        streams.close();
    }
}
