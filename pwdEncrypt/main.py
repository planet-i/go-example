import base64
import hashlib
import rsa


def md5_ecc(str):
    md = hashlib.md5()  # 创建md5对象
    md.update(str.encode(encoding='utf-8'))
    return md.hexdigest()

def rsa_ecc(str):
    msg = str.encode()
    # msg = b'aaaaaaaaa-4444'
    public_key_str = '-----BEGIN PUBLIC KEY-----\nxxxxxxxx\n-----END PUBLIC KEY-----\n'
    # (pub_key, priv_key) = rsa.newkeys(256)
    public_key = rsa.PublicKey.load_pkcs1_openssl_pem(public_key_str.encode())
    ctxt = rsa.encrypt(msg, public_key)
    print(base64.b64encode(ctxt).decode())
    return (base64.b64encode(ctxt).decode())

def login(user_name=name, password=pwd):
    """
    用户登录，默认admin用户
    :param user_name:
    :param password:
    :return:
    """
    url = "/api/as/user/login"
    data = {"loginId": user_name, "password": rsa_ecc(md5_ecc(password)), "from": 0, "type": 0}
    status_code, response_data = req.post(url=url, params=data)
    if response_data.get('code') == 200 and response_data.get("result"):
        return response_data['result']['token'], response_data['result']['info']['root_ids']
    else:
        raise Exception("Error! User login fail!!!")