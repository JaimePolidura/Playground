package es.jaime.curator;

import lombok.RequiredArgsConstructor;
import org.apache.curator.framework.CuratorFramework;
import org.apache.curator.framework.recipes.leader.LeaderSelector;
import org.apache.curator.framework.recipes.leader.LeaderSelectorListener;
import org.apache.curator.framework.state.ConnectionState;

import java.io.Closeable;
import java.io.IOException;
import java.util.UUID;

@RequiredArgsConstructor
public final class CuratorLeaderSelector implements LeaderSelectorListener, Closeable {
    private final LeaderSelector leaderSelector;
    private final UUID nodeId;

    public void runForMaster() {
        this.leaderSelector.setId(this.nodeId.toString());
        this.leaderSelector.start();
    }

    @Override
    public void takeLeadership(CuratorFramework curatorFramework) throws Exception {

    }

    @Override
    public void stateChanged(CuratorFramework curatorFramework, ConnectionState connectionState) {

    }

    @Override
    public void close() throws IOException {

    }
}
