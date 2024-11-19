package io.poc.articles_ms;

import java.util.List;
import java.util.Map;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import io.dapr.client.DaprClient;
import io.dapr.exceptions.DaprException;


@RestController
@CrossOrigin
@RequestMapping("/api/order")
public class OrdersController {

    @Autowired
    private  DaprClient daprClient;

    @Autowired
    private OrderRepo orderRepo;

    @Value("${dapr.secretstore}")
    private String secretStore;

    @Value("${dapr.pubsubname}")
    private String pubSubName;

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
              Map<String, String> pubsub = daprClient.getSecret(secretStore, pubSubName).block();
              System.out.println(pubsub.get("pubSubName"));
              daprClient.publishEvent(pubsub.get("pubSubName"), "orders", order).block();
            return ResponseEntity.status(HttpStatus.OK).body("Order Placed successfully");
        } catch (DaprException e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("Error sending message: " + e.getMessage());
        }
    }






}
