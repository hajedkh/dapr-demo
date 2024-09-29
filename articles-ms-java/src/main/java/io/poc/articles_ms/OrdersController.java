package io.poc.articles_ms;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import io.dapr.client.DaprClient;
import io.dapr.exceptions.DaprException;

@RestController
public class OrdersController {

    @Autowired
    private  DaprClient daprClient;



    @PostMapping("/send")
    public ResponseEntity<String> sendOrderMessage(@RequestBody Order order) {
        try {
              daprClient.publishEvent("rabbitmq-pubsub", "orders", order).block();
            return ResponseEntity.status(HttpStatus.OK).body("Order Placed successfully");
        } catch (DaprException e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("Error sending message: " + e.getMessage());
        }
    }

    @PostMapping("/test")
    public void test(){
        System.out.println("message recieved");
    }

}
