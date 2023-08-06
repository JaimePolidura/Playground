package es.jaime.leader;

import org.apache.zookeeper.KeeperException;
import org.apache.zookeeper.ZooKeeper;
import org.apache.zookeeper.data.Stat;

import java.util.UUID;

public final class LeaderUtils {
    public static boolean isLeaderTaken(ZooKeeper zooKeeper, UUID nodeId) {
        while (true) {
            try {
                Stat stat = new Stat();
                byte[] response = zooKeeper.getData("/leader", false, stat);

                return UUID.nameUUIDFromBytes(response).equals(nodeId);
            } catch (KeeperException.NodeExistsException e) {
                return false;
            } catch (Exception ignored) {}
        }
    }
}
