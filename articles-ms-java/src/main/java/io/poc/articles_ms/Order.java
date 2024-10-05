package io.poc.articles_ms;


import jakarta.persistence.*;

import java.util.ArrayList;
import java.util.List;

/**
 * just for test reason this is supposed to be a DTO not an entity
 */
@Entity
@Table(name = "orders")
public class Order {

    @Id
    private String id;


    private List<String> articleIds = new ArrayList<>();
    private int quantity;

    // Constructors, Getters, Setters

    public Order() {}

    public Order(String orderId,List<String> articleIds, int quantity) {
        this.articleIds = articleIds;
        this.id = orderId;
        this.quantity = quantity;
    }

    public String getOrderId() {
        return id;
    }

    public void setOrderId(String orderId) {
        this.id = orderId;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }



    public int getQuantity() {
        return quantity;
    }

    public void setQuantity(int quantity) {
        this.quantity = quantity;
    }

    public List<String> getArticleIds() {
        return articleIds;
    }

    public void setArticleIds(List<String> articleIds) {
        this.articleIds = articleIds;
    }
}

