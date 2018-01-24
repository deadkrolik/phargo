<?php
$p = new Phar('alias_md5.phar');

$p->addFromString('data.txt', 'DATA');
$p->setSignatureAlgorithm(\Phar::MD5);
$p->setAlias('ALIAS');
