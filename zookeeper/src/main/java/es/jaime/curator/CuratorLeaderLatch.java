package es.jaime.curator;

import lombok.RequiredArgsConstructor;
import lombok.SneakyThrows;
import org.apache.curator.framework.CuratorFramework;
import org.apache.curator.framework.api.CuratorEvent;
import org.apache.curator.framework.api.CuratorListener;
import org.apache.curator.framework.recipes.leader.LeaderLatch;
import org.apache.curator.framework.recipes.leader.LeaderLatchListener;

import java.io.Closeable;
import java.io.IOException;

@RequiredArgsConstructor
public final class CuratorLeaderLatch implements Closeable, LeaderLatchListener {
    private final CuratorFramework client;
    private final LeaderLatch leaderLatch;

    public void runForMaster() throws Exception {
        client.getCuratorListenable().addListener(new MasterListener());
        leaderLatch.addListener(this);
        leaderLatch.start();
    }

    @Override
    @SneakyThrows
    public void isLeader() {
        //When I get the leadership
        leaderLatch.start();
    }

    @Override
    public void notLeader() {
        //When I lost the leadership
    }

    @Override
    public void close() throws IOException {
        leaderLatch.close();
    }

    private static class MasterListener implements CuratorListener {
        @Override
        public void eventReceived(CuratorFramework curatorFramework, CuratorEvent event) throws Exception {
            switch (event.getType()) {
                case CHILDREN -> {
                    System.out.println("Successfully got a list of assignments: " + event.getChildren().size() + " tasks");
                    //deleteAssignment(event.getPath() + "/" + task);
                    //deleteAssignment(event.getPath());
                    //assignTasks(event.getChildren());
                }
                case CREATE -> {
                    //deleteTask(event.getPath().substring(event.getPath().lastIndexOf('-') + 1));
                }
            }
        }
    }
}
