package es.jaime.helloworld;

import lombok.SneakyThrows;
import org.apache.zookeeper.WatchedEvent;
import org.apache.zookeeper.ZooKeeper;

public class HellWorld {
    @SneakyThrows
    public static void main(String[] args) {
        try (ZooKeeper zooKeeper = new ZooKeeper("127.0.0.1:2181", 5000, HellWorld::onWatchEvent)) {
            System.out.println("Connected to zookeeper!");
        }
    }

    private static void onWatchEvent(WatchedEvent watchedEvent) {
        System.out.println("Received watch event: " + watchedEvent.getType());
    }
}
