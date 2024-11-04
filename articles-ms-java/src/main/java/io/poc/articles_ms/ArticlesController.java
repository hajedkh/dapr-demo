package io.poc.articles_ms;


import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@CrossOrigin
@RequestMapping("/api/article")
public class ArticlesController {

    @Autowired
    ArticleRepo articleRepo;

    @GetMapping("/getAll")
    public List<Article> getAll(){

        return articleRepo.findAll();
    }

    @PostMapping("/add")
    public ResponseEntity<String> addArticles(@RequestBody List<Article> articles){
        try {
            articleRepo.saveAll(articles);
            return ResponseEntity.status(HttpStatus.OK).body("articles added successfully !");
        } catch (Exception e){
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                    .body("Error sending message: " + e.getMessage());
        }

    }
}
