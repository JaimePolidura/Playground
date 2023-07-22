package es.jaime.trees;

import lombok.AllArgsConstructor;
import lombok.Getter;

import java.util.function.Consumer;

public class BinarySearchTree<T extends Comparable<T>> implements Tree<T> {
    private Node<T> rootNode;
    private int size;

    @Override
    public boolean insert(T data) {
        if(this.rootNode == null){
            this.rootNode = new Node<>(data, null, null);
            this.size++;
            return true;
        }

        Node<T> actualNode = this.rootNode;

        return this.addRecursive(data, this.rootNode);
    }

    private boolean addRecursive(T data, Node<T> parentNode) {
        if(parentNode.data.equals(data))
            return false;

        boolean actualDataBigger = !parentNode.isBiggerThan(data);

        if(actualDataBigger)
            return this.addChildNode(parentNode.right, newNode -> parentNode.right = newNode, data, parentNode);
        else
            return this.addChildNode(parentNode.left, newNode -> parentNode.left = newNode, data, parentNode);
    }

    private boolean addChildNode(Node<T> childNodeToTraverse , Consumer<Node<T>> newNodeParentSetter, T newData, Node<T> parentNode) {
        if(childNodeToTraverse == null){
            Node<T> newNode = new Node<>(newData, null, null);
            newNodeParentSetter.accept(newNode);
            this.size++;

            return true;
        }else{
            return this.addRecursive(newData, childNodeToTraverse);
        }
    }

    @Override
    public boolean search(T data) {
        return !this.isEmpty() && this.searchRecursive(this.rootNode, this.rootNode, data) != null;
    }

    @Override
    public boolean delete(T data) {
        if(this.isEmpty())
            return false;

        NodeSearchResult<T> nodeSearchResult = searchRecursive(this.rootNode, this.rootNode, data);
        if(nodeSearchResult != null && nodeSearchResult.isEmpty())
            return false;

        deleteNode(nodeSearchResult.getNode(), nodeSearchResult.getParentNode());

        return true;
    }

    private void deleteNode(Node<T> nodeToDelete, Node<T> parentNodeToDelete) {
        if(nodeToDelete.isLeaf() && parentNodeToDelete.left == nodeToDelete)
            parentNodeToDelete.left = null;

        if(nodeToDelete.isLeaf() && parentNodeToDelete.right == nodeToDelete)
            parentNodeToDelete.right = null;

        if(nodeToDelete.onlyHasOneChildLeft())
            parentNodeToDelete.left = nodeToDelete.left;

        if(nodeToDelete.onlyHasOneChildRight())
            parentNodeToDelete.right = nodeToDelete.right;

        if(nodeToDelete.isInnerNode() && nodeToDelete == this.rootNode){
            NodeSearchResult<T> searchResultNodeToReplace = this.findEdgeNodeInSubtree(this.rootNode, false);
            Node<T> nodeToReplace = searchResultNodeToReplace.getNode();

            Node<T> nodeLeftToRoot = this.rootNode.left != nodeToReplace ? this.rootNode.left : null;
            Node<T> nodeRightToRoot = this.rootNode.right != nodeToReplace ? this.rootNode.right : null;

            deleteNode(nodeToReplace, searchResultNodeToReplace.getParentNode());

            nodeToReplace.left = nodeLeftToRoot;
            nodeToReplace.right = nodeRightToRoot;
            this.rootNode = nodeToReplace;
        }

        if(nodeToDelete.isInnerNode() && nodeToDelete != this.rootNode){
            NodeSearchResult<T> searchResultNodeToReplace = parentNodeToDelete.right == nodeToDelete ?
                this.findEdgeNodeInSubtree(nodeToDelete.right, false) :
                this.findEdgeNodeInSubtree(nodeToDelete.left, true);
            Node<T> nodeToReplace = searchResultNodeToReplace.getNode();

            deleteNode(nodeToReplace, searchResultNodeToReplace.getParentNode());

            nodeToReplace.left = nodeToDelete.left;
            nodeToReplace.right = nodeToDelete.right;

            if(parentNodeToDelete.right == nodeToDelete)
                parentNodeToDelete.right = nodeToDelete;
            else
                parentNodeToDelete.left = nodeToDelete;
        }

        this.size--;
    }


    private NodeSearchResult<T> findEdgeNodeInSubtree(Node<T> subtreeRootNode, boolean goLeft) {
        Node<T> parentToMinimunNode = subtreeRootNode;
        Node<T> minimunNode = goLeft ? subtreeRootNode.left : subtreeRootNode.right;

        while (minimunNode != null && (goLeft ? minimunNode.hasLeft() : minimunNode.hasRight())){
            parentToMinimunNode = minimunNode;
            minimunNode = goLeft ? minimunNode.right : minimunNode.left;
        }

        return new NodeSearchResult<>(minimunNode, parentToMinimunNode);
    }

    @Override
    public int size() {
        return this.size;
    }

    @Override
    public void clear() {
        this.rootNode = null;
        this.size = 0;
    }

    private NodeSearchResult<T> searchRecursive(Node<T> parentNode, Node<T> node, T dataToSearch) {
        if(node.getData().equals(dataToSearch))
            return new NodeSearchResult<>(node, parentNode);

        return node.isBiggerThan(dataToSearch) ?
            (node.hasLeft() ? searchRecursive(node, node.left, dataToSearch) : null) :
            (node.hasRight() ? searchRecursive(node, node.right, dataToSearch) : null);
    }

    @AllArgsConstructor
    private static class NodeSearchResult<T extends Comparable<T>> {
        @Getter private final Node<T> node;
        @Getter private final Node<T> parentNode;

        public boolean isEmpty() {
            return this.node == null;
        }
    }

    @AllArgsConstructor
    private static class Node<T extends Comparable<T>> {
        @Getter private final T data;
        @Getter private Node<T> left;
        @Getter private Node<T> right;

        public boolean isBiggerThan(T data) {
            return this.data.compareTo(data) > 0;
        }

        public boolean hasRight() {
            return this.right != null;
        }

        public boolean hasLeft() {
            return this.left != null;
        }

        public boolean onlyHasOneChildLeft() {
            return this.left != null && this.right == null;
        }

        public boolean onlyHasOneChildRight() {
            return this.left == null && this.right != null;
        }

        public boolean hasTwoChildren() {
            return this.left != null && this.right != null;
        }

        public boolean isInnerNode() {
            return this.left != null && this.right != null;
        }

        public boolean isLeaf() {
            return this.left == null && this.right == null;
        }
    }
}
