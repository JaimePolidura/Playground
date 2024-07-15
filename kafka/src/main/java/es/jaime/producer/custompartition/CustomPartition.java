package es.jaime.producer.custompartition;

import org.apache.kafka.clients.producer.Partitioner;
import org.apache.kafka.common.Cluster;

import java.util.Map;

public final class CustomPartition implements Partitioner {
    @Override
    public int partition(
            String topic,
            Object key,
            byte[] keyBytes,
            Object value,
            byte[] valueBytes,
            Cluster cluster
    ) {
        int numPartitions = cluster.partitionsForTopic(topic).size();

        return (int) (Math.random() * numPartitions);
    }

    @Override
    public void close() {
        //?
    }

    @Override
    public void configure(Map<String, ?> map) {
        //Vacio
    }
}
