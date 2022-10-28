<?php
$method = "GET";

$url = "http://pdxf.malhyuk.info:8888/api/v2/adversaries?include=name&include=adversary_id";

$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $url);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 10);
curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);

$headers = [
    "KEY: ADMIN123",
    "accept: application/json",
    "Content-Type: application/json; charset=utf-8"
];
curl_setopt($ch, CURLOPT_HEADER, true);
curl_setopt($ch, CURLOPT_HTTPHEADER, $headers);

$response = curl_exec($ch);
if ($response == false) {
    $error = curl_error($ch);
    echo $error;
}
curl_close($ch);
$code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
echo $code;

$header_size = curl_getinfo($ch, CURLINFO_HEADER_SIZE);
$header = substr($response, 0, $header_size);
$body = substr($response, $header_size);

$v = json_decode($body,true);
header('Content-Type: application/json');
echo json_encode($v, JSON_PRETTY_PRINT);
// echo $body;
?>