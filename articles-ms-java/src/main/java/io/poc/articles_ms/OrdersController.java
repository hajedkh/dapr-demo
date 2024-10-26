package io.poc.articles_ms;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import io.dapr.client.DaprClient;
import io.dapr.exceptions.DaprException;

import java.util.List;

@RestController
@CrossOrigin
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


    @GetMapping("/get")
    public List<Order> getOrders(){
         return orderRepo.findAll();
    }



    @PostMapping("/send")
    public ResponseEntity<String> sendOrderMessage(@RequestBody Order order) {
        try {
              daprClient.publishEvent("redis-pubsub", "orders", order).block();
            return ResponseEntity.status(HttpStatus.OK).body("Order Placed successfully");
        } catch (DaprException e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("Error sending message: " + e.getMessage());
        }
    }






}
