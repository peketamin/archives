<?php
function accumulator($x){
    $f = function($y) use ($x){
        return $x + $y;
    };
    return $f;
}

$a = accumulator(10);
foreach(range(0, 5) as $i) {
    echo $a($i)."\n";
}
