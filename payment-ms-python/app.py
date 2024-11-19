import os
import time

from flask import Flask, jsonify, request
from dapr.clients import DaprClient
import json

PUBSUB_NAME = 'redis-pubsub'
SECRETSTORE_NAME = 'secretstore'
TOPIC_NAME = 'payments'

app = Flask("Payment-MS")

base_url = os.getenv('BASE_URL', 'http://localhost') + ':' + os.getenv(
                    'DAPR_HTTP_PORT', '3500')

class Order:
    def __init__(self, orderId=None, articleIds=None, quantity=0):
        self.orderId = orderId
        self.articleIds = articleIds
        self.quantity = quantity

    def __str__(self):
        return f"Order {self.orderId} ({self.quantity}x {self.articleIds})"

@app.route('/hello')
def hello():
    return 'Hello, World!'

@app.route('/pay', methods=['POST'])
def pay():
    with DaprClient() as client:
        secret = client.get_secret(store_name=SECRETSTORE_NAME, key="payment-key")
    print("starting the shit!!")
    order = Order(request.json['orderId'], request.json['articleIds'], request.json['quantity'])
    print(f"Received payment query for {order} and using secret key {secret.secret}", flush=True)
    time.sleep(5)
    print(f"Payment succeeded for {order}", flush=True)
    with DaprClient() as d:
        result = d.publish_event(
            pubsub_name=PUBSUB_NAME,
            topic_name=TOPIC_NAME,
            data=json.dumps({'id': order.orderId, 'success': True}),
            data_content_type='application/json',
        )
    return jsonify({'id': order.orderId, 'success': True}), 200

if __name__ == '__main__':
    app.run(port=5001)
