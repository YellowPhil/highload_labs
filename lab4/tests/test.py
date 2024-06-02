import requests
from time import sleep, time
from random import choices, choice
from string import printable
from hashlib import sha256
from concurrent.futures import ThreadPoolExecutor

URL = "http://localhost:8080"
WORDS = []
ROUNDS = 500
MAX_THREADS = 8


def random_string(n: int=15):
    return ''.join(choices(printable, k=n))


def submit_hash(hash: str):
    requests.post(f"{URL}/add_task", data = {"hash": hash})


def generate_hash() -> str:
    word = choice(WORDS)
    hash = sha256(word.encode()).digest()
    return hash.hex()


def generator():
    for _ in range(ROUNDS):
        submit_hash(generate_hash())

def generate_load(threads: int):
    pool = ThreadPoolExecutor(max_workers=threads)
    for _ in range(threads - 1):
        pool.submit(generator)
    pool.submit(queue_len)
    pool.shutdown(wait = True)


def queue_len():
    start_time = time()
    while True:
        with open("results.txt", "a") as f:
            f.write(f"{time() - start_time} {requests.get(f"{URL}/queue_len").text}\n")
        sleep(0.5)


if __name__ == '__main__':
    WORDS = open("./wordlist.txt").read().split("\n")
    for _ in range(200):
        WORDS.append(random_string())
    generate_load(MAX_THREADS)
    # queue_len()
