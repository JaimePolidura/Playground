package es.jaime.workers;

import lombok.RequiredArgsConstructor;
import org.apache.zookeeper.*;
import org.apache.zookeeper.data.Stat;

import java.util.concurrent.ConcurrentHashMap;

@RequiredArgsConstructor
public final class Client {
    private final ZooKeeper zookeeper;

    public void submit(String taskCommand) {
        zookeeper.create("/tasks/task-",
                taskCommand.getBytes(),
                ZooDefs.Ids.OPEN_ACL_UNSAFE,
                CreateMode.PERSISTENT_SEQUENTIAL,
                this::onTaskCreatedCallback,
                taskCommand);
    }

    private void onTaskCreatedCallback(int resultCode, String path, Object ctx, String name, Stat stat) {
        String taskCommand = (String) ctx;

        switch (KeeperException.Code.get(resultCode)) {
            case CONNECTIONLOSS -> submit(taskCommand);
            case OK -> watchStatus("/status/" + name.replace("/tasks/", ""), taskCommand);
        }
    }

    ConcurrentHashMap<String, Object> ctxMap = new ConcurrentHashMap<String, Object>();

    private void watchStatus(String watchPath, String taskCommand) {
        ctxMap.put(watchPath, taskCommand);
        zookeeper.exists(watchPath,
                this::taskStatusWatcher,
                this::existsTaskCallback,
                taskCommand);
    }

    private void existsTaskCallback(int resultCode, String path, Object ctx, Stat stat) {
        String taskCommand = (String) ctx;

        switch (KeeperException.Code.get(resultCode)) {
            case CONNECTIONLOSS -> watchStatus(path, taskCommand);
            case OK -> zookeeper.getData(path, false, this::getTaskDataCallback, null);
        }
    }

    private void taskStatusWatcher(WatchedEvent event) {
        if (event.getType() == Watcher.Event.EventType.NodeCreated) {
            zookeeper.getData(event.getPath(), false, this::getTaskDataCallback, null);
        }
    }

    private void getTaskDataCallback(int i, String s, Object o, byte[] bytes, Stat stat) {
        System.out.println("BOOM MINESHAFT!");
    }
}
