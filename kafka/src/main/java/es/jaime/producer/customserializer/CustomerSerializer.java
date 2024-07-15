package es.jaime.producer.customserializer;

import org.apache.kafka.common.serialization.Serializer;

import java.nio.ByteBuffer;
import java.nio.charset.StandardCharsets;

public final class CustomerSerializer implements Serializer<Customer> {
    @Override
    public byte[] serialize(String topic, Customer customer) {
        int stringSize = customer.getName().length();
        ByteBuffer buffer = ByteBuffer.allocate(4 + 4 + stringSize);
        buffer.putInt(customer.getCustomerId());
        buffer.putInt(stringSize);
        buffer.put(customer.getName().getBytes(StandardCharsets.UTF_8));

        return buffer.array();
    }
}
