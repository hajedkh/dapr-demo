package io.poc.articles_ms;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;


/**
 * Just For test reasons to not be used
 */
@Repository
public interface OrderRepo extends JpaRepository<Order, String> {

}
