package es.jaime.helloworld;

import lombok.SneakyThrows;
import org.apache.zookeeper.*;

import java.util.Arrays;
import java.util.List;

public class HellWorld {
    @SneakyThrows
    public static void main(String[] args) {
        try (ZooKeeper zookeeper = new ZooKeeper("127.0.0.1:2181", 5000, HellWorld::onWatchEvent)) {
            System.out.println("Connected to zookeeper!");

            List<OpResult> results = zookeeper.multi(List.of(
                    Op.delete("/nada", -1),
                    Op.setData("/caca", "nada".getBytes(), -1)
            ));
        }
    }

    private static void onWatchEvent(WatchedEvent watchedEvent) {
        System.out.println("Received watch event: " + watchedEvent.getType());
    }
}
