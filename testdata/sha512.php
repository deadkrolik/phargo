<?php
$p = new Phar('sha512.phar');

$p['FILE'] = 'FDATA';
$p->setSignatureAlgorithm(\Phar::SHA512);
