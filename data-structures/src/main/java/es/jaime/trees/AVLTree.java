package es.jaime.trees;

import lombok.AllArgsConstructor;
import lombok.Getter;

public final class AVLTree<T extends Comparable<T>> implements Tree<T> {
    private Node<T> root;
    private int size;

    @Override
    public boolean search(T data) {
        Node<T> node = this.root;

        while (node != null) {
            if (data.equals(node.data))
                return true;
            else if (data.compareTo(node.data) < 0)
                node = node.left;
            else
                node = node.right;
        }

        return false;
    }

    @Override
    public boolean delete(T data) {
        this.deleteNodeRecursive(this.root, data);

        return true;
    }

    private Node<T> deleteNodeRecursive(Node<T> last, T toRemove) {
        if (last == null) {
            return last;
        } else if (last.data.compareTo(toRemove) > 0) {
            last.left = deleteNodeRecursive(last.left, toRemove);
        } else if (last.data.compareTo(toRemove) < 0) {
            last.right = deleteNodeRecursive(last.right, toRemove);
        } else {
            this.size--;

            if (last.left == null || last.right == null) {
                last = (last.left == null) ? last.right : last.left;
            } else {
                this.size++;

                Node<T> mostLeftChild = this.mostLeftChild(last.right);
                last.data = mostLeftChild.data;
                last.right = deleteNodeRecursive(last.right, last.data);
            }
        }

        if (last != null)
            last = rebalance(last);

        return last;
    }

    private Node<T> mostLeftChild(Node<T> node) {
        while (node.left != null)
            node = node.left;

        return node;
    }

    @Override
    public void clear() {
        this.root = null;
        this.size = 0;
    }

    @Override
    public int size() {
        return this.size;
    }

    @Override
    public boolean insert(T data) {
        Node<T> newNode = new Node<>(data, null, null, -1);

        if(this.root == null)
            this.root = newNode;
        else
            this.insertRecursive(this.root, data);

        this.size++;

        return true;
    }

    private Node<T> insertRecursive(Node<T> last, T data) {
        if (last == null)
            return new Node(data, null, null, -1);
        else if (last.data.compareTo(data) > 0)
            last.left = insertRecursive(last.left, data);
        else if (last.data.compareTo(data) < 0)
            last.right = insertRecursive(last.right, data);

        return rebalance(last);
    }

    private Node<T> rebalance(Node<T> node) {
        this.updateHeight(node);

        int balanceFactor = this.getHeightFactor(node);

        if (balanceFactor < -1) { //Left heavy
            if (this.getHeightFactor(node.left) > 0)
                node.left = rotateLeft(node.left);

            node = rotateRight(node);
        }

        if (balanceFactor > 1) { //Right heavy
            if (this.getHeightFactor(node.right) < 0)
                node.right = rotateRight(node.right);
            node = rotateLeft(node);
        }

        return node;
    }

    private Node<T> rotateRight(Node<T> node) {
        Node<T> leftChild = node.left;

        this.updateRootReferenceIfNeccesary(node, leftChild);

        node.left = leftChild.right;
        leftChild.right = node;

        this.updateHeight(node);
        this.updateHeight(leftChild);

        return leftChild;
    }

    private Node<T> rotateLeft(Node node) {
        Node<T> rightChild = node.right;

        node.right = rightChild.left;
        rightChild.left = node;

        this.updateRootReferenceIfNeccesary(node, rightChild);

        this.updateHeight(node);
        this.updateHeight(rightChild);

        return rightChild;
    }

    private void updateRootReferenceIfNeccesary(Node<T> oldReference, Node<T> newReference) {
        if(this.root == oldReference)
            this.root = newReference;
    }

    private int getHeightFactor(Node<T> node) {
        return node == null ? 0 : getHeight(node.getRight()) - getHeight(node.getLeft());
    }

    private void updateHeight(Node<T> node) {
        int leftChildHeight = this.getHeight(node.left);
        int rightChildHeight = this.getHeight(node.right);

        node.height = Math.max(leftChildHeight, rightChildHeight) + 1;
    }

    private int getHeight(Node<T> node) {
        return node != null ? node.getHeight() : -1;
    }

    @AllArgsConstructor
    private static class Node<T extends Comparable<T>> {
        @Getter private T data;
        @Getter private Node<T> left;
        @Getter private Node<T> right;
        @Getter private int height;
    }
}
