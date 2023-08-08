package es.jaime.curator;

import org.apache.curator.framework.CuratorFramework;
import org.apache.curator.framework.CuratorFrameworkFactory;
import org.apache.curator.framework.recipes.cache.PathChildrenCacheEvent;
import org.apache.curator.framework.recipes.cache.PathChildrenCacheListener;
import org.apache.curator.framework.recipes.leader.LeaderSelector;
import org.apache.curator.framework.recipes.leader.LeaderSelectorListener;
import org.apache.curator.framework.recipes.locks.InterProcessMultiLock;
import org.apache.curator.framework.state.ConnectionState;
import org.apache.curator.retry.RetryNTimes;
import org.apache.zookeeper.CreateMode;

import java.util.UUID;

public final class Main {
    public static void main(String[] args) throws Exception {
        try(var curator = CuratorFrameworkFactory.newClient("127.0.0.1:2181", new RetryNTimes(10, 2000))) {
            curator.getCuratorListenable().addListener((curatorFramework, curatorEvent) -> {
                switch (curatorEvent.getType()) {

                }
            });

            curator.create()
                    .withMode(CreateMode.EPHEMERAL)
                    .forPath("/nodes");

            curator.create()
                    .withProtection() //Seq node will have a prefix with a unique id, so that if the zookeeper retries the operation, it won't create two seq nodes
                    .withMode(CreateMode.EPHEMERAL_SEQUENTIAL)
                    .inBackground()
                    .forPath("/nodes", "datos".getBytes());


            PathChildrenCacheListener pathChildrenCacheListener = (curatorFramework, pathChildrenCacheEvent) -> {
                if (pathChildrenCacheEvent.getType() == PathChildrenCacheEvent.Type.CHILD_REMOVED) {

                }
            };

            LeaderSelector leaderSelector = new LeaderSelector(curator, "/master", null);
        }
    }
}
