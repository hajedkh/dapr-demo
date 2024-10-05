package io.poc.articles_ms;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import io.dapr.client.DaprClient;
import io.dapr.exceptions.DaprException;

@RestController
@RequestMapping("/api/order")
public class OrdersController {

    @Autowired
    private  DaprClient daprClient;

    @Autowired
    private OrderRepo orderRepo;



    /** Just For test reasons do not use it ,
    /*  the front needs to cache the order created and send it directly threw the /send endpoint just below
     */
    @PostMapping("/add")
    public ResponseEntity<String> addOrder(@RequestBody Order order){
        try {
            orderRepo.save(order);
            return ResponseEntity.status(HttpStatus.OK).body("order added successfully !");
        } catch (Exception e){
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("Error sending message: " + e.getMessage());
        }
    }







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



}
