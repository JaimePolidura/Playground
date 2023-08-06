package es.jaime.leader;

import es.jaime.utils.ManualFuture;
import lombok.Getter;
import lombok.RequiredArgsConstructor;
import org.apache.zookeeper.*;
import org.apache.zookeeper.data.Stat;

import java.util.UUID;
import java.util.concurrent.Future;

@RequiredArgsConstructor
public final class SimpleAsyncNode {
    private final ZooKeeper zooKeeper;
    private final UUID nodeId;

    @Getter private boolean isLeader;

    public Future<Boolean> tryToBeLeader() {
        ManualFuture<Boolean> future = ManualFuture.ofDefault(false);

        this.zooKeeper.create("/leader",
                this.nodeId.toString().getBytes(),
                ZooDefs.Ids.OPEN_ACL_UNSAFE,
                CreateMode.EPHEMERAL,
                (resultCode, s, o, s1, stat) -> this.onLeaderCreateCallback(resultCode, s, o, s1, stat, future),
                null);

        return future;
    }

    private void onLeaderCreateCallback(int resultCode, String s, Object o, String s1, Stat stat, ManualFuture<Boolean> leaderFuture) {
        switch (KeeperException.Code.get(resultCode)) {
            case CONNECTIONLOSS -> this.isLeader = LeaderUtils.isLeaderTaken(zooKeeper, nodeId);
            case OK -> this.isLeader = true;
            default -> this.isLeader = false;
        }

        leaderFuture.complete(this.isLeader);
    }

    public static void main(String[] args) throws Exception {
        try (var zookeeper = new ZooKeeper("127.0.0.1:2181", 5000, e -> {})) {
            SimpleAsyncNode node1 = new SimpleAsyncNode(zookeeper, UUID.randomUUID());
            SimpleAsyncNode node2 = new SimpleAsyncNode(zookeeper, UUID.randomUUID());

            System.out.println("Should be true:" + node1.tryToBeLeader().get() + " " + node1.isLeader());
            System.out.println("Should be false:" + node2.tryToBeLeader().get() + " " + node2.isLeader());
        }
    }

}
