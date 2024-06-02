import requests
import redis
import json
from flask import Flask
from base64 import b64decode, b64encode
from os import getenv
from time import time

RUNS = 100

app = Flask(__name__)
redis_client = redis.Redis(
    host = 'localhost',
    port = int(
        getenv('REDIS_PORT', '6379')
    )
)


def get_dog_picture(width: int):
    cached = redis_client.get(str(width))
    if cached is not None:
        return b64decode(cached.decode())

    dog_picture = requests.get(f"https://place.dog/300/{width}", stream=True)

    redis_client.set(str(width), b64encode(dog_picture.content), ex=60)
    return dog_picture


def get_agregated_dog_picture(width: int):
# Представим себе тут НЕВЕРОЯТНУЮ СИСТЕМУ АГРЕГАЦИИ для картинок с собаками
    if width % 6 == 0:
        width = 200
    return get_dog_picture(width)



def dry_run():
    start_time = time()
    for i in range(200, 200 + RUNS):
        get_dog_picture(i)

    return time() - start_time


def agregated_run():
    start_time = time()
    for i in range(200, 200 + RUNS):
        get_agregated_dog_picture(i)

    return time() - start_time


def warmed_up_run():
    for i in range(200, 200 + RUNS):
        get_dog_picture(i)

    start_time = time()
    for i in range(200, 200 + RUNS):
        get_dog_picture(i)

    return time() - start_time


if __name__ == '__main__':
    print(f"Dry run: {dry_run()}")
    redis_client.flushall()

    print(f"Argegated run: {dry_run()}")
    redis_client.flushall()

    print(f"Warmed up run: {warmed_up_run()}")


