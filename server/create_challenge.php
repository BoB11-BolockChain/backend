<?php

// <!--문제 번호 / OS / 이미지 번호 / 문제 제목 / 문제 설명 / 점수 / 공격 개수 / 공격 플래그-->

// CREATE TABLE training (
    // num integer not null 
    // primary key autoincrement, 
    // os varchar(10) not null, 
    // img_num int not null, 
    // title varchar(50) not null, 
    // scenario varchar(500) not null, 
    // score int not null, 
    // atk_cnt int not null, 
    // atk_flag varchar(30) // for problem process
    // )

$DB = new SQLite3('../front/db');
if($DB->lastErrorCode() == 0){
	echo "Database connection succeed!";
}
else {
	echo "Database connection failed";
    echo $DB->lastErrorMsg();
}

$os = $_POST["OSdata"];
$server = $_POST["ServerType"];
$_attack = $_POST["AttackData"];
$attack1 = $_attack["attack1"];
$attack2 = $_attack["attack2"];
$attack3 = $_attack["attack3"];

$DB->exec(`insert into training values ({$os},{$},{$},{$},{$},{$},{$})`);


?>