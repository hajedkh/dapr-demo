package io.poc.articles_ms;





public class Order {
    private String orderId;
    private String articleId;
    private int quantity;

    // Constructors, Getters, Setters

    public Order() {}

    public Order(String orderId, String articleId, int quantity) {
        this.orderId = orderId;
        this.articleId = articleId;
        this.quantity = quantity;
    }

    public String getOrderId() {
        return orderId;
    }

    public void setOrderId(String orderId) {
        this.orderId = orderId;
    }

    public String getArticleId() {
        return articleId;
    }

    public void setArticleId(String articleId) {
        this.articleId = articleId;
    }

    public int getQuantity() {
        return quantity;
    }

    public void setQuantity(int quantity) {
        this.quantity = quantity;
    }
}

