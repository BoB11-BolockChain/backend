<?php
$os = $_POST["payload"];

if (!isset($_POST["payload"])) {
    exit("error:invalid form");
}

// run docker container
// wait container ready or schedule(??) operation
// send docker access point to challenger

// deploy agent - exe_name: name
$res_ability = shell_exec("python3 Upload_Attack.py name & echo \"agent deploy\"");

//create operation
$method = "POST";

// Required schema fields are as follows: "name", "adversary.adversary_id", "planner.planner_id", and "source.id"
$data = array(
    "name"=>"{$name}",
    "adversary"=>array(
        "adversary_id"=>"6c038d82-8dfc-4490-b321-e9c20e142c59"
    ),
    "planner"=>array(
        "id"=>"aaa7c857-37a0-4c4a-85f7-4e9f7f30e31a" //default stockpile atomic planner
    ),
    "source"=>array(
        "id"=>"ed32b9c3-9593-4c33-b0db-e2007315096b" //default basic source
    )
);

$url = "http://pdxf.malhyuk.info:8888/api/v2/operations";

$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $url);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
curl_setopt($ch, CURLOPT_CONNECTTIMEOUT, 10);
curl_setopt($ch, CURLOPT_SSL_VERIFYPEER, false);
//curl_setopt($ch, CURLOPT_SSLVERSION, 3); // https
curl_setopt($ch, CURLOPT_CUSTOMREQUEST, $method);
// $post_field_string=http_build_query($data,'','&');
curl_setopt($ch,CURLOPT_POSTFIELDS,json_encode($data));


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
$header_size = curl_getinfo($ch, CURLINFO_HEADER_SIZE);
$header = substr($response, 0, $header_size);
$body = substr($response, $header_size);
$v = json_decode($body,true);
header('Content-Type: application/json');
echo $code . "\n";
echo json_encode($v, JSON_PRETTY_PRINT);
?>