<?php
$fname = 'bad_hash.phar';

@unlink($fname);
$p = new Phar($fname);

$p->addFromString('111', '222');
$p->setSignatureAlgorithm(\Phar::MD5);

$f = fopen($fname, 'r+');
fseek($f, 50);
fwrite($f, '---------');
fclose($f);
