package es.jaime.workers;

import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import org.apache.zookeeper.ZooKeeper;

@RequiredArgsConstructor
public final class Admin {
    private final ZooKeeper zookeeper;
    
    @SneakyThrows
    public void showState() {
        System.out.println("Workers:");
        for (String w: zookeeper.getChildren("/workers", false)) {
            byte[] data = zookeeper.getData("/workers/" + w, false, null);
            String state = new String(data);
            System.out.println("\t" + w + ": " + state);
        }
        System.out.println("Tasks:");
        for (String t: zookeeper.getChildren("/tasks", false)) {
            System.out.println("\t" + t);
        }
    }
}
