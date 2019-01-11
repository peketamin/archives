<?php
$terms = array(
    'namespace' => '名前空間',
    'requirements' => '要件',
    'array' => '配列 (キーワードとして使用されているところでは array)',
    'parameters' => 'パラメータ',
    'static' => '静的',
    'returns' => '返り値',
    'example' => '例',
    'browser' => 'ブラウザ',
    'request' => 'リクエスト (HTTP リクエストの場合)',
    'default' => 'デフォルト',
    'note' => '注意',
    'folder' => 'フォルダ',
    'usage' => '使用法',
    'web' => 'Web',
    'we' => '私たち',
    'requirement' => '要件',
    'Fuel' => 'Fuel',
    'user' => 'ユーザ',
    'webserver' => 'Web サーバ',
    'encoding' => 'エンコーディング',
    'introduction' => 'はじめに',
    'prefix' => 'プレフィックス',
    'permission' => 'パーミッション',
    'action' => 'アクション',
    'role' => 'ロール',
);

function echo_hilight_term($str, $term) {
    $str = trim($str);
    $words = explode($str, ' ');
    $offset = 0;
    while ($offset < strlen($str)) {
        $pos = stripos($str, $term, $offset);
        if ($pos === false) {
            echo substr($str, $offset);
            break;
        }
        else {
            echo substr($str, $offset, $pos);
            $hilight_word = substr($str, $pos, strlen($term));
            $color_start = "\033[31m";
            $color_terminate = "\033[37m";
            echo $color_start.$hilight_word.$color_terminate;
            $offset = $pos + strlen($term);
        }
    }
}
function find_terms ($str, $terms, $line_number) {
    foreach($terms as $orig_term => $trans_term) {
        $pat = preg_quote($orig_term, '#');
        $found = preg_match("#\b$pat\b#i", $str, $matches);
        if ($found) {
            echo 'LN: '.$line_number.PHP_EOL;
            // echo trim($str).PHP_EOL;
            echo_hilight_term($str, $orig_term);
            echo "\n$orig_term => $trans_term\n\n";
        }
    }
}
function main($argc, $argv) {
    global $terms;

    if ($argc < 2) {
        echo "次の引数が不足しています: 読み込み対象ファイル名\n";
        return;
    }

    $filename = $argv[1];
    if ( ! file_exists($filename)) {
        echo 'ファイルが存在しません: '.$filename.PHP_EOL;
        return;
    }

    $lines = file($filename);
    $in_body = false;
    $after_footer = false;
    $ln = 0;
    foreach ($lines as $line) {
        $ln++;
        if (strpos($line, '<body>') !== false) {
            $in_body = true;
        }
        if ( ! $in_body) {
            continue;
        }
        if (strpos($line, '<footer>') !== false) {
            $after_footer = true;
        }
        if ($after_footer) {
            continue;
        }
        find_terms($line, $terms, $ln);
    }
}

main($argc, $argv);
