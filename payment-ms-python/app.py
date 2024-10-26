import os
import time

from flask import Flask, jsonify, request
from dapr.clients import DaprClient
import json

PUBSUB_NAME = 'redis-pubsub'
TOPIC_NAME = 'payments'

app = Flask("Payment-MS")

base_url = os.getenv('BASE_URL', 'http://localhost') + ':' + os.getenv(
                    'DAPR_HTTP_PORT', '3500')

class Order:
    def __init__(self, order_id=None, article_ids=None, quantity=0):
        self.order_id = order_id
        self.article_ids = article_ids
        self.quantity = quantity

    def __str__(self):
        return f"Order {self.order_id} ({self.quantity}x {self.article_ids})"

@app.route('/hello')
def hello():
    return 'Hello, World!'

@app.route('/pay', methods=['POST'])
def pay():
    print("starting the shit!!")
    order = Order(request.json['order_id'], request.json['article_ids'], request.json['quantity'])
    print(f"Received payment query for {order}", flush=True)
    time.sleep(5)
    print(f"Payment succeeded for {order}", flush=True)
    with DaprClient() as d:
        result = d.publish_event(
            pubsub_name=PUBSUB_NAME,
            topic_name=TOPIC_NAME,
            data=json.dumps({'id': order.order_id, 'success': True}),
            data_content_type='application/json',
        )
    return jsonify({'id': order.order_id, 'success': True}), 200

if __name__ == '__main__':
    app.run(port=5001)
