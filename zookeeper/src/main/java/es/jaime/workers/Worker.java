package es.jaime.workers;

import io.netty.util.internal.ConcurrentSet;
import lombok.SneakyThrows;
import org.apache.zookeeper.*;
import org.apache.zookeeper.data.Stat;

import java.util.*;

public final class Worker {
    private final ZooKeeper zookeeper;
    private final UUID nodeId;

    private State state;

    //Worker
    private Set<UUID> onGoingTaskIds;

    //Master
    private List<UUID> workersIds;

    public Worker(ZooKeeper zooKeeper, UUID workerId) {
        this.onGoingTaskIds = Collections.synchronizedSet(new HashSet<>());
        this.workersIds = Collections.synchronizedList(new ArrayList<>());
        this.zookeeper = zooKeeper;
        this.state = State.IDLE;
        this.nodeId = workerId;
    }

    public void initialize() {
        this.zookeeper.exists("/master",
                this::masterExistsWatch,
                this::masterExistsCallback,
                null);
    }

    private void setState(State state) {
        this.zookeeper.setData(getWorkerPath(),
                state.toString().getBytes(),
                -1,
                this::onSetStateCallback,
                state);
    }

    private void masterExistsCallback(int resultCode, String s, Object o, Stat stat) {
        switch (KeeperException.Code.get(resultCode)) {
            case CONNECTIONLOSS -> initialize();
            case NODEEXISTS -> startWorker();
            case OK -> startMaster();
        }
    }

    private void startWorker() {
        this.state = State.IDLE;

        this.zookeeper.create(getWorkerPath(),
                State.IDLE.toString().getBytes(),
                ZooDefs.Ids.OPEN_ACL_UNSAFE,
                CreateMode.EPHEMERAL,
                this::onRegisterCallback,
                null);
    }

    private void startMaster() {
        this.state = State.MASTER;
        this.deleteSelfWorker();

        handleWorkers();
        handleTasks();
    }

    private void handleWorkers() {
        this.zookeeper.getChildren("/workers",
                this::onWorkersChangedWatch,
                this::onWorkerChildren,
                null);
    }

    private void onWorkerChildren(int resultCode, String s, Object o, List<String> workersIdsString, Stat stat) {
        switch (KeeperException.Code.get(resultCode))  {
            case CONNECTIONLOSS -> startMaster();
            case OK -> this.workersIds = workersIdsString.stream()
                    .map(UUID::fromString)
                    .toList();
        }
    }

    private void onWorkersChangedWatch(WatchedEvent event) {
        this.handleWorkers();
    }

    private void handleTasks() {
        this.zookeeper.getChildren("/tasks",
                this::onTasksChangedWatch,
                this::onTasksCallback,
                null);
    }

    private void onTasksChangedWatch(WatchedEvent event) {
        if (event.getType() == Watcher.Event.EventType.NodeChildrenChanged) {
            handleTasks();
        }
    }

    private void onTasksCallback(int resultCode, String s, Object o, List<String> unassignedTasksIdStrings, Stat stat) {
        switch (KeeperException.Code.get(resultCode)) {
            case CONNECTIONLOSS -> this.handleTasks();
            case OK -> assignUnassignedTasks(unassignedTasksIdStrings.stream().map(UUID::fromString).toList());
        }
    }

    private void assignUnassignedTasks(List<UUID> unassignedTasksIds) {
        for (UUID unassignedTaskId : unassignedTasksIds) {
            getTaskDataAndAssignTask(unassignedTaskId);
        }
    }

    private void getTaskDataAndAssignTask(UUID taskId) {
        zookeeper.getData("/tasks/" + taskId.toString(),
                false,
                this::onGetTaskData,
                taskId);
    }

    private void onGetTaskData(int resultCode, String s, Object o, byte[] bytes, Stat stat) {
        String taskCommand = new String(bytes);
        UUID taskId = (UUID) o;

        switch (KeeperException.Code.get(resultCode)) {
            case CONNECTIONLOSS -> getTaskDataAndAssignTask((UUID) o);
            case OK -> assignTask(taskId, taskCommand);
        }
    }

    private void assignTask(UUID unassignedTaskId, String command) {
        UUID workerId = workersIds.get((int) (Math.random() * workersIds.size()));
        String assignmentPath = "/assign/" + workerId.toString() + "/" + unassignedTaskId.toString();

        createTasksAssigment(unassignedTaskId, assignmentPath, command);
    }

    @SneakyThrows
    private void createTasksAssigment(UUID taskId, String assignmentPath, String command) {
        zookeeper.multi(List.of(
                Op.create(assignmentPath, command.getBytes(), ZooDefs.Ids.OPEN_ACL_UNSAFE, CreateMode.PERSISTENT),
                Op.delete("/tasks/" + taskId.toString(), -1)
        ));
    }

    private void deleteSelfWorker() {
        while (true) {
            try {
                this.zookeeper.delete("/workers/" + this.nodeId.toString(), -1);
            } catch (KeeperException.ConnectionLossException e) {
                //continue;
            } catch (Exception e){
                return;
            }
        }
    }

    private void masterExistsWatch(WatchedEvent event) {
        if (event.getType() == Watcher.Event.EventType.NodeDeleted) {
            initialize();
        }
    }

    private void onSetStateCallback(int resultCode, String s, Object o, Stat stat) {
        switch (KeeperException.Code.get(resultCode)) {
            case CONNECTIONLOSS -> setState((State) o);
            case OK -> this.state = (State) o;
        }
    }

    private void onRegisterCallback(int resultCode, String s, Object o, String s1, Stat stat) {
        switch (KeeperException.Code.get(resultCode)) {
            case CONNECTIONLOSS -> startWorker();
            case OK -> getWorkerAssignedTasks();
            case NODEEXISTS -> System.out.printf("Worker %s already exists\n", nodeId.toString());
        }
    }

    private void getWorkerAssignedTasks() {
        zookeeper.getChildren("/assign/" +
                this.nodeId,
                this::onNewAssignedTaskWatch,
                this::onTaskAssignmentsGetCallback,
                null);
    }

    private void onTaskAssignmentsGetCallback(int resultCode, String s, Object o, List<String> assignedTasksStringId, Stat stat) {
        var assignedTasksId = assignedTasksStringId.stream().map(UUID::fromString).toList();

        switch (KeeperException.Code.get(resultCode)) {
            case CONNECTIONLOSS -> getWorkerAssignedTasks();
            case OK -> {
                for (UUID assignedTaskId : assignedTasksId) {
                    if (onGoingTaskIds.contains(assignedTaskId)) {
                        onGoingTaskIds.add(assignedTaskId);
                        //Execute task
                        //Delete task from assigment & tasks & onGoingTaskIds
                    }
                }
            }
        }
    }

    private void onNewAssignedTaskWatch(WatchedEvent event) {
        if (event.getType() == Watcher.Event.EventType.NodeChildrenChanged) {
            this.getWorkerAssignedTasks();
        }
    }

    private String getWorkerPath() {
        return "/workers/" + this.nodeId.toString();
    }

    public enum State {
        MASTER, IDLE
    }
}
