import requests
import hashlib
import json
import client_config as cfg
from User import User


if __name__ == '__main__':
    cur_user = User()
    while True:
        command = input()
        print(command)
