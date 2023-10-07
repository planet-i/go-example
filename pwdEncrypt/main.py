from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import padding
from cryptography.hazmat.primitives import hashes
import base64

# 公钥字符串
user_public_key = """
-----BEGIN PUBLIC KEY-----
xxxxxxx
-----END PUBLIC KEY-----
"""

# 要加密的密码
password = "Aa123456789+".encode()  # 需要将密码编码为字节串

# 将公钥字符串解析为公钥对象
public_key = serialization.load_pem_public_key(user_public_key.encode(), backend=default_backend())

# 加密数据
encrypted_data = public_key.encrypt(
    password,
    padding.OAEP(
        mgf=padding.MGF1(algorithm=hashes.SHA256()),
        algorithm=hashes.SHA256(),
        label=None
    )
)

# 将加密后的数据转换为 Base64 字符串
base64_encrypted_data = base64.b64encode(encrypted_data).decode()
print(base64_encrypted_data)
