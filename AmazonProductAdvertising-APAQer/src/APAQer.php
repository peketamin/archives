<?php
namespace Peketamin\AmazonProductAdvertising;
use Peketamin\AmazonProductAdvertising\APAQer;
use Peketamin\AmazonProductAdvertising\APAQerException;
use Peketamin\AmazonProductAdvertising\APAQerConfig;
use Peketamin\AmazonProductAdvertising\APAQerConfigException;

class APAQerConfigException extends \Exception {};
class APAQerConfig {
    public $jsonFile = "";
    public $requestParams = [
        'AWSAccessKeyId' => 'xxxxxxxxxxxxxxxxxxxx',
        'AWSSecretKeyId' => 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx',
        'AssociateTag' => 'xxxxxxxxxxxxxxxxxxxx',

        'Operation' => 'ItemSearch',
        'ResponseGroup' => 'Large,Images',
        'SearchIndex' => 'Books',
        'Sort' => 'daterank',
        'Power' => 'language:English',
        'BrowseNode' => '2275256051',
        'MaximumPrice' => '0',
        'Condition' => 'New',
        'ItemPage' => '1',
        'Timestamp' => '',
    ];

    public function __construct($jsonFile) {
        $this->jsonFile = $jsonFile;
        $dataArr = static::_readFile($this->jsonFile);

        // filtering and set user settings in json config to instance settings.
        foreach ($this->requestParams as $key => $value) {
            if (array_key_exists($key, $dataArr)) {
                $this->requestParams[$key] = $dataArr[$key];
            }
        }

        if ($this->requestParams['SearchIndex'] == 'KindleStore') {
            unset($this->requestParams['Power']);
        }

    }

    public function resetItemPage() {
        $this->requestParams['ItemPage'] = 1;
    }

    private static function _readFile($jsonFile) {
        if ( ! file_exists($jsonFile)) {
            throw new APAQerConfigException('config file could not find.');
        }

        $jsonText = file_get_contents($jsonFile);
        $phpArr = json_decode($jsonText, true);
        if ( ! ($phpArr) or ! is_array($phpArr)) {
            throw new APAQerConfigException('json decoding failed.');
        }

        return $phpArr;
    }
}


class APAQerException extends \Exception {};
class APAQerRequestException extends APAQerException {};
class APAQerIsValidException extends APAQerException {};
class APAQer {
    // to reference easy.
    protected $docurl = array(
        'http://docs.aws.amazon.com/AWSECommerceService/2013-08-01/DG/AnatomyOfaRESTRequest.html',
    );

    public $requestURL = "";
    public $config = null;

    const DATE_FORMAT = 'Y-m-d\TH:i:m\Z';
    protected static $defaultTimezone = "";
    protected static $amazonHost = 'webservices.amazon.co.jp';
    protected static $baseUrl = "";
    protected $requestBase = "";
    protected $signature = "";
    public $totalPages = "";

    public function __construct(APAQerConfig $config) {
        // Amazon API needs UTC timezone
        static::$defaultTimezone = date_default_timezone_get();
        date_default_timezone_set('UTC');

        static::$baseUrl = 'http://'.static::$amazonHost.'/onca/xml?';

        $this->config = $config;
        $this->config->requestParams['Timestamp'] = date(self::DATE_FORMAT);
        ksort($this->config->requestParams);
        date_default_timezone_set(static::$defaultTimezone);
    }

    private function _makeRequestBase() {
        $this->requestBase = http_build_query($this->config->requestParams, null, '&', PHP_QUERY_RFC3986);
        return $this->requestBase;
    }

    private function _makeSignature() {
        $request_with_header = "GET\n".static::$amazonHost."\n/onca/xml\n{$this->requestBase}";
        $hash =  base64_encode(hash_hmac('sha256', $request_with_header, $this->config->requestParams['AWSSecretKeyId'], true));
        $hash_enc = urlencode($hash);
        $this->signature = $hash_enc;
        return $this->signature;
    }

    public function makeRequestFull() {
        $this->_makeRequestBase();
        $this->_makeSignature();
        $this->requestURL = static::$baseUrl.$this->requestBase.'&Signature='.$this->signature;

        return $this->requestURL;
    }

    public function fetch() {
        $xml = @file_get_contents($this->makeRequestFull());
        if (strpos($http_response_header[0], '200') == false) {
            return false;
        }
        $responseObjTree = @simplexml_load_string($xml);

        $xmlRequestErrors = libxml_get_errors();
        if ($xmlRequestErrors) {
            //throw new APAQerRequestException("request failed.");
            return false;
        }

        if (isset($responseObjTree->Items->Request->IsValid) and $responseObjTree->Items->Request->IsValid != "True") {
            //throw new APAQerIsValidException("response: request is not valid.");
            return false;
        }

        if ( ! isset($responseObjTree->Items->TotalPages)) {
            return false;
        }

        $this->totalPages = $responseObjTree->Items->TotalPages;

        return $xml;
    }

    public static function getItemNodes($xml) {
        $responseObjTree = @simplexml_load_string($xml);
        if (isset($responseObjTree->Items->Item)) {
            return $responseObjTree->Items->Item;
        }
        return [];
    }

    //public function __destruct() {
    //    date_default_timezone_set(static::$defaultTimezone);
    //}
}
