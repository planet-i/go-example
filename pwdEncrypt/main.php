const USER_PUBLIC_KEY = "-----BEGIN PUBLIC KEY-----
xxxxxxx
-----END PUBLIC KEY-----";

$password = md5("Aa123456789+");
// 转换为公钥资源
$publicKey = openssl_pkey_get_public(USER_PUBLIC_KEY);
// 加密数据
openssl_public_encrypt($password, $encryptedData, $publicKey);
// 将加密后的数据转换为 Base64 字符串
$base64EncryptedData = base64_encode($encryptedData);
var_dump($base64EncryptedData);