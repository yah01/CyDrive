import client_config as cfg
import requests
import ast
import json
import hashlib
from User import User


class BaseCmd:

    def __init__(self):
        pass


    def get_cmd_map(self):
        cmd_map = {
            'login': self.login,
            'list': self.query,
            'exit': self.exit,
            'exit_account': self.exit_account,
            'download': self.download,

            '登录': self.login,
            '查询': self.query,
            '退出': self.exit,
            '注销': self.exit_account,
            '下载': self.download,
        }
        return cmd_map


    def execute(self, args=''):
        if '，' in args:
            args = args.replace('，', ',')
        args = args.split()

        cmd_map = self.get_cmd_map()
        if args[0] not in cmd_map:
            return False, 'No such Command!'
        cmd_func = cmd_map[args[0]]
        parser_args = []
        for item in args[1:]:
            try:
                item = ast.literal_eval(item)
                parser_args.append(item)
            except ValueError as err:
                parser_args.append(item)
        try:
            status, msg = cmd_func(*parser_args)
            return status, msg
        except Exception as err:
            return False, err


    def login(self, username=None, password=None):
        def judge(txt):
            txt = json.loads(txt)
            if txt['status'] == 0:
                return True
            return False

        if username is None and password is None:
            username = 'test'
            password = 'testCyDrive'

        password = password.encode()
        psw_hashed = hashlib.sha256(hashlib.md5(password).digest()).digest()
        psw_str = ''
        for item in psw_hashed:
            psw_str = psw_str + str(int(item))

        login_response = requests.post(cfg.URLS['login'], data={'username': username, 'password': psw_str})
        global user_cookie
        # global cur_user
        user_cookie = login_response.cookies

        login_dict = json.loads(login_response.text)
        if judge(login_response.text):
            if username == 'test':
                return True, '测试账号登陆成功！' + login_dict['message']
            return True, '登陆成功！' + login_dict['message']
            # cur_user.
        return False, '登陆失败！' + login_dict['message']


    def query(self, path=None):
        if path is None:
            path = ''
        global user_cookie
        lists_response = requests.get(cfg.URLS['list'] + '?' + 'path=' + path, cookies=user_cookie)

        list_res = json.loads(lists_response.text)
        msg = '查询成功！'
        if list_res['status'] != 0:
            msg = '查询失败！'
        if list_res['status'] == 0:
            for item in list_res['data']:
                print(item)
        return list_res['status'] == 0, msg


    def exit(self):
        print('886')
        exit(0)


    def exit_account(self):
        global user_cookie
        try:
            user_cookie = None
        except Exception as err:
            return False, '注销失败，错误信息：\n' + str(err)
        return True, '注销成功！'


    def download(self, file_path=None):
        if file_path is None:
            file_path = ''
        global user_cookie
        download_response = requests.get(cfg.URLS['download'] + '?' + 'filepath=' + file_path, cookies=user_cookie)
        print(download_response.text)
        response_dict = json.loads(download_response)
        status = response_dict['status']
        if status == 2:
            return False, '下载失败，只能下载单个文件！'
        msg = '下载成功！'
        if status != 0:
            msg = '下载失败！'
        return status == 0, msg


if __name__ == '__main__':
    user_cookie = None
    cmd = BaseCmd()
    while True:
        command = input()
        status, msg = cmd.execute(command)
        print(status, msg)