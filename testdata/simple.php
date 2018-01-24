<?php
$p = new Phar('simple.phar');

$p->addFromString('1.txt', 'ASDF');
$p->addFromString('index.php', 'ZXCV');
$p->setSignatureAlgorithm(\Phar::SHA1);

$p->setMetadata([
    'a' => 123,
]);
