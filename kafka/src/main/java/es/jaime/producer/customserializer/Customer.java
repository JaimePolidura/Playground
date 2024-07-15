package es.jaime.producer.customserializer;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;

@Getter
@Builder
@AllArgsConstructor
public final class Customer {
    private final int customerId;
    private final String name;

}
