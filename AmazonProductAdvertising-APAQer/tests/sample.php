<?php
error_reporting(E_ALL);
ini_set("display_errors", '1');

require_once('../vendor/autoload.php');

use Peketamin\AmazonProductAdvertising\APAQer;
use Peketamin\AmazonProductAdvertising\APAQerConfig;

// ----------------------------------------

$ITEM_COUNTER = 0;
$PAGE_COUNTER = 0;

// for some display HTML elements
$echoHeader_1 = function () {
    return '<h1>Sample using Amazon API</h1>'.PHP_EOL;
};

$echoHeaderPageCount = function ($pageCounter) {
    return sprintf('<hr><h2>Page: %d</h2>', $pageCounter).PHP_EOL;
};

// ### Display list of items as HTML
$echoHTML = function($xml, $PAGE_COUNTER, &$ITEM_COUNTER) {

    // ### Display item counter
    $echoItemCounter = function($localItemCounter, $PAGE_COUNTER) {
        $_counter = ($PAGE_COUNTER - 1) * 10 + $localItemCounter;
        return '<em>&#35;'.$_counter.'</em>';
    };


    $resultHTML = '<meta charset="UTF-8">'.PHP_EOL;
    $localItemCounter = 0;

    $resultHTML .= "<!--{$xml}-->";  // display raw XML Strings as HTML comments.
    $resultHTML .= "<ul>\n";

    $items = APAQer::getItemNodes($xml);
    foreach ($items as $item) {
        // count items up
        ++$ITEM_COUNTER;
        ++$localItemCounter;

        $resultHTML .= '<!--'.print_r($item, true).'-->';

        $resultHTML .=
        '<li class=box>
            <table>
                <tr>
                    <td>
                        <i>'.$echoItemCounter($localItemCounter, $PAGE_COUNTER).'</i>
                    </td>
                    <td>
                        <img src="'.$item->MediumImage->URL.'">
                    </td>
                    <td>
                        <a href="'.$item->DetailPageURL.'">'.$item->ItemAttributes->Title.'</a>
                    </td>
                </tr>
            </table>
        </li>';
    }
    $resultHTML .= "</ul>\n";

    return $resultHTML;
};

// CSS
// ---
$css =<<<__CSS__
<style>
.box {
    width: 50%;
    border: 1px solid gray;
}
</style>
__CSS__;


// Setup
// =================
$amazonConfig = new APAQerConfig('APAQerConfig.json');
$amazonApi = new APAQer($amazonConfig);
$amazonApi->config->requestParams['BrowseNode'] = '2293143051';


// Display Main page
// =================
echo $css;

$MAIN_HTML = '';
$mainHTML = function($addingHTML) use (&$MAIN_HTML) {
    $MAIN_HTML .= $addingHTML;
    echo $addingHTML;
};

$mainHTML($echoHeader_1());


// Page 1
// ------
$queryResponse = $amazonApi->fetch();
$PAGE_COUNTER++;

$mainHTML($echoHeaderPageCount($PAGE_COUNTER));
$mainHTML($echoHTML($queryResponse, $PAGE_COUNTER, $ITEM_COUNTER));


// display more pages.
if ($amazonApi->totalPages > 1) {
    $maxPages = ($amazonApi->totalPages < 10) ? $amazonApi->totalPages : 10;
    // $maxPages = 3; // test
    for($PAGE_COUNTER = 2; $PAGE_COUNTER <= $maxPages; $PAGE_COUNTER++) {
        $amazonApi->config->requestParams['ItemPage'] = $PAGE_COUNTER;
        //print_r($amazonApi->config->requestParams);
        $queryResponse = $amazonApi->fetch();

        $mainHTML($echoHeaderPageCount($PAGE_COUNTER));
        $mainHTML($echoHTML($queryResponse, $PAGE_COUNTER, $ITEM_COUNTER));
    }
}

// ### Summary
$mainHTML(sprintf('<p>Count of All item is: %s</p>'.PHP_EOL, $ITEM_COUNTER));

//echo $MAIN_HTML;
exit;

