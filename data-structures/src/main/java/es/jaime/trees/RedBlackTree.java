package es.jaime.trees;

import lombok.AllArgsConstructor;
import lombok.Getter;

public final class RedBlackTree<T extends Comparable<T>> implements Tree<T> {
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
    public boolean insert(T data) {
        Node<T> parentNode = this.findNodeToInsert(data);
        Node<T> newNode = new Node(data, null, null, null, NodeColor.RED);

        newNode.parent = parentNode;
        if (parentNode == null)
            this.root = newNode;
        else if (data.compareTo(parentNode.data) < 0)
            parentNode.left = newNode;
        else
            parentNode.right = newNode;

        this.size++;

        fixRedBlackPropertiesAfterInsert(newNode);

        return true;
    }

    private void fixRedBlackPropertiesAfterInsert(Node<T> child) {
        Node<T> parent = child.parent;
        if(parent == null || parent.color == NodeColor.BLACK)
            return;

        Node<T> grandparent = parent.parent;

        if(grandparent == null){
            parent.color = NodeColor.BLACK;
            return;
        }

        Node<T> uncle = this.getUncle(parent);
        if (uncle != null && uncle.color == NodeColor.RED) {
            parent.color = NodeColor.BLACK;
            grandparent.color = NodeColor.RED;
            uncle.color = NodeColor.BLACK;

            this.fixRedBlackPropertiesAfterInsert(grandparent);
            return;
        }

        if(grandparent.left == parent){
            if(parent.right == child) { //Triangulo
                this.rotateLeft(parent);
                parent = child;
            }

            this.rotateRight(grandparent); //Linea

        } else if(grandparent.right == parent){
            if(parent.left == child){ //Triangulo
                this.rotateRight(parent);
                parent = child;
            }

            this.rotateLeft(grandparent); //Linea
        }

        parent.color = NodeColor.BLACK; //Pasa a ser padre del nuevo subtrianuglo.
        grandparent.color = NodeColor.RED; //Pasa a ser hermano del nodo insertado
    }

    private Node getUncle(Node<T> parent) {
        Node grandparent = parent.parent;
        if (grandparent.left == parent)
            return grandparent.right;
        else if (grandparent.right == parent)
            return grandparent.left;
        else
            return null;
    }

    private Node<T> findNodeToInsert(T data) {
        Node<T> node = this.root;
        Node<T> parent = null;

        while (node != null){
            parent = node;

            if (data.compareTo(node.data) < 0)
                node = node.left;
            else if (data.compareTo(node.data) > 0)
                node = node.right;
            else
                return node;
        }

        return node == null ? parent : node;
    }

    @Override
    public boolean delete(T data) {
        Node<T> node = this.findNodeToInsert(data);

        if (node == null) {
            return false;
        }

        Node<T> movedUpNode;
        NodeColor deletedNodeColor;

        if (node.left == null || node.right == null) {
            movedUpNode = deleteNodeWithZeroOrOneChild(node);
            deletedNodeColor = node.color;
        } else {
            Node<T> inOrderSuccessor = findMinimum(node.right);

            node.data = inOrderSuccessor.data;
            movedUpNode = deleteNodeWithZeroOrOneChild(inOrderSuccessor);
            deletedNodeColor = inOrderSuccessor.color;
        }

        if (deletedNodeColor == NodeColor.BLACK) {
            fixRedBlackPropertiesAfterDelete(movedUpNode);

            if (movedUpNode.getClass() == NilNode.class)
                replaceParentsChild(movedUpNode.parent, movedUpNode, null);
        }

        this.size--;

        return true;
    }

    private void fixRedBlackPropertiesAfterDelete(Node node) {
        if (node == root)
            return;

        Node<T> sibling = this.getSibling(node);

        if(sibling == null)
            return;

        if (sibling.color == NodeColor.RED){
            handleRedSibling(node, sibling);
            sibling = getSibling(node);
        }

        if(sibling == null)
            return;

        if (isBlack(sibling.left) && isBlack(sibling.right)) {
            sibling.color = NodeColor.BLACK;

            if (node.parent.color == NodeColor.RED)
                node.parent.color = NodeColor.BLACK;
            else
                fixRedBlackPropertiesAfterDelete(node.parent);
        } else {
            handleBlackSiblingWithAtLeastOneRedChild(node, sibling);
        }
    }

    private Node<T> getSibling(Node<T> node) {
        Node parent = node.parent;
        if (node == parent.left)
            return parent.right;
        else if (node == parent.right)
            return parent.left;
        else
            return null;
    }

    private void handleBlackSiblingWithAtLeastOneRedChild(Node node, Node sibling) {
        boolean nodeIsLeftChild = node == node.parent.left;

        if (nodeIsLeftChild && isBlack(sibling.right)) {
            sibling.left.color = NodeColor.BLACK;
            sibling.color = NodeColor.RED;
            rotateRight(sibling);
            sibling = node.parent.right;
        } else if (!nodeIsLeftChild && isBlack(sibling.left)) {
            sibling.right.color = NodeColor.BLACK;
            sibling.color = NodeColor.RED;
            rotateLeft(sibling);
            sibling = node.parent.left;
        }

        sibling.color = node.parent.color;
        node.parent.color = NodeColor.BLACK;
        if (nodeIsLeftChild) {
            sibling.right.color = NodeColor.BLACK;
            rotateLeft(node.parent);
        } else {
            sibling.left.color = NodeColor.BLACK;
            rotateRight(node.parent);
        }
    }

    private void handleRedSibling(Node<T> node, Node<T> sibling) {
        sibling.color = NodeColor.BLACK;
        node.parent.color = NodeColor.RED;

        if (node == node.parent.left)
            rotateLeft(node.parent);
        else
            rotateRight(node.parent);
    }

    private boolean isBlack(Node node) {
        return node == null || node.color == NodeColor.BLACK;
    }

    private Node<T> findMinimum(Node<T> node) {
        while (node.left != null)
            node = node.left;

        return node;
    }

    private Node<T> deleteNodeWithZeroOrOneChild(Node<T> node) {
        if (node.left != null) {
            this.replaceParentsChild(node.parent, node, node.left);
            return node.left;
        }

        if (node.right != null) {
            this.replaceParentsChild(node.parent, node, node.right);
            return node.right;
        }

        Node<T> newChild = node.color == NodeColor.BLACK ? new NilNode() : null;
        this.replaceParentsChild(node.parent, node, newChild);

        return newChild;
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

    private void rotateLeft(Node<T> node) {
        Node<T> parent = node.parent;
        Node<T> rightChild = node.right;

        node.right = rightChild.left;
        if (rightChild.left != null)
            rightChild.left.parent = node;

        rightChild.left = node;
        node.parent = rightChild;

        replaceParentsChild(parent, node, rightChild);
    }

    private void rotateRight(Node<T> node) {
        Node<T> parent = node.parent;
        Node<T> leftChild = node.left;

        node.left = leftChild.right;
        if (leftChild.right != null)
            leftChild.right.parent = node;

        leftChild.right = node;
        node.parent = leftChild;

        replaceParentsChild(parent, node, leftChild);
    }

    private void replaceParentsChild(Node<T> parent, Node<T> oldChild, Node<T> newChild) {
        if (parent == null)
            this.root = newChild;
        else if (parent.left == oldChild)
            parent.left = newChild;
        else if (parent.right == oldChild)
            parent.right = newChild;

        if (newChild != null)
            newChild.parent = parent;
    }

    private static class NilNode<T extends Comparable<T>> extends Node<T> {
        public NilNode() {
            super(null, null, null, null, NodeColor.BLACK);
        }
    }

    @AllArgsConstructor
    private static class Node<T extends Comparable<T>> {
        @Getter private T data;
        @Getter private Node<T> left;
        @Getter private Node<T> right;
        @Getter private Node<T> parent;
        @Getter private NodeColor color;

        public boolean isLeaf() {
            return this.left == null && this.right == null;
        }

        public boolean hasOneChild() {
            return (this.left == null && this.right != null) || (this.left != null && this.right == null);
        }
    }

    private enum NodeColor {
        RED, BLACK;
    }
}
