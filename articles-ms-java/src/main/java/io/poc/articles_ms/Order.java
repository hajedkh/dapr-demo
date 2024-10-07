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
    private String order_id;


    private List<String> article_ids = new ArrayList<>();
    private int quantity;

    // Constructors, Getters, Setters

    public Order() {}

    public Order(String order_id, List<String> article_ids, int quantity) {
        this.article_ids = article_ids;
        this.order_id = order_id;
        this.quantity = quantity;
    }


    public String getOrder_id() {
        return order_id;
    }

    public void setOrder_id(String id) {
        this.order_id = id;
    }



    public int getQuantity() {
        return quantity;
    }

    public void setQuantity(int quantity) {
        this.quantity = quantity;
    }

    public List<String> getArticle_ids() {
        return article_ids;
    }

    public void setArticle_ids(List<String> articleIds) {
        this.article_ids = articleIds;
    }
}

