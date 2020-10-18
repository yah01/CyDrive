import requests
import hashlib
import json
import client_config as cfg


class User:
    def __init__(self, user_id=-1, username='', password='', work_dir='', root_dir='', remote_dir=''):
        self.user_id = user_id
        self.username = username
        self.password = password
        self.work_dir = work_dir
        self.root_dir = root_dir
        self.remote_dir = remote_dir

    def print(self):
        print(self.username)
        print(self.password)

    def instanciate(self):
        pass

    def user_unmarshal(self, src_json):
        new_instance = json.loads(src_json, )
        return new_instance

