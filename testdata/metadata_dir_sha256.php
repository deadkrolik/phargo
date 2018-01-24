<?php
$p = new Phar('metadata_dir_sha256.phar');

$p['FILE'] = 'FDATA';
$p['FILE']->setMetadata(['v' => 'x']);
$p['/DIR1/FILE1'] = 'D1_DATA11';
$p['/DIR1/FILE2'] = 'D1_DATA12';
$p['/DIR2/FILE1'] = 'D1_DATA21';
$p['/DIR2/FILE1']->setMetadata(['z' => 'cc']);
$p->setSignatureAlgorithm(\Phar::SHA256);
