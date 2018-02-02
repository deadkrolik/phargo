<?php
$p = new Phar('gz.phar');

$p['ABCD'] = 'DATADATADATADATA';
$p['ABCD']->compress(\Phar::GZ);
