package es.jaime.leader;

import lombok.Getter;
import lombok.RequiredArgsConstructor;
import org.apache.zookeeper.CreateMode;
import org.apache.zookeeper.KeeperException;
import org.apache.zookeeper.ZooDefs;
import org.apache.zookeeper.ZooKeeper;

import java.util.UUID;

@RequiredArgsConstructor
public final class SimpleSyncNode {
    private final ZooKeeper zooKeeper;
    private final UUID nodeId;

    @Getter private boolean isLeader;

    public boolean tryToBeLeader() {
        do {
            try {
                this.zooKeeper.create(
                        "/leader",
                        this.nodeId.toString().getBytes(),
                        ZooDefs.Ids.OPEN_ACL_UNSAFE,
                        CreateMode.EPHEMERAL
                );
                this.isLeader = true;
                return true;
            } catch (KeeperException.NodeExistsException e) {
                this.isLeader = false;
                return false;
            } catch (Exception e) {
                e.printStackTrace();
            }
        } while(!(this.isLeader = LeaderUtils.isLeaderTaken(zooKeeper, nodeId)));

        return true;
    }

    public static void main(String[] args) throws Exception {
        try (var zookeeper = new ZooKeeper("127.0.0.1:2181", 5000, e -> {})) {
            SimpleSyncNode node1 = new SimpleSyncNode(zookeeper, UUID.randomUUID());
            SimpleSyncNode node2 = new SimpleSyncNode(zookeeper, UUID.randomUUID());

            System.out.println("Should be true:" + node1.tryToBeLeader() + " " + node1.isLeader());
            System.out.println("Should be false:" + node2.tryToBeLeader() + " " + node2.isLeader());
        }
    }
}
