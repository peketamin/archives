AmazonProductAdvertising-APAQer
===============================

Sample Library using AmazonProductAdvertising API


Overview
--------

Just a sample code and library of AmazonProductAdvertising API.
This is a prototype just for for me.


## Description

see a screenshot in *Demo* section below...


## Demo

![screenshot](https://raw.github.com/wiki/peketamin/AmazonProductAdvertising-APAQer/images/screenshot.png?raw=true)


## VS.

Maybe many better lib on github or on the internet...


## Requirement

- PHP >= 5.4.3
- [Composer](https://getcomposer.org/)


## Usage

### 1

Make JSON config file.
Skelton file is in `src/`.
Copy and customize it.


### 2

coding as you like.

```php
<?php
require_once('../vendor/autoload.php');

use Peketamin\AmazonProductAdvertising\APAQer;
use Peketamin\AmazonProductAdvertising\APAQerConfig;

// Read config from json file.
$amazonConfig = new APAQerConfig('APAQerConfig.json');
// Make instance of this library.
$amazonApi = new APAQer($amazonConfig);

// Set broweNode:
// * en: [Browse Node IDs - Product Advertising API](http://docs.aws.amazon.com/AWSECommerceService/latest/DG/BrowseNodeIDs.html)
// * ja: [Amazon アソシエイト（アフィリエイト） - ヘルプ](https://affiliate.amazon.co.jp/gp/associates/help/t100)
$amazonApi->config->requestParams['BrowseNode'] = '2293143051';

// Get Request result
$queryResponse = $amazonApi->request();
```

See more in tests/sample.php .


## Install

Assuming you know about [Composer](https://getcomposer.org/).

1. Add a line to composer.json like this.
   ```json:composer.json
       ...
   	"require": {
   		...
   		"peketamin/AmazonProductAdvertising-APAQer": "@dev"
   	},
   	...
   ```
2. `composer install` or `composer update`


## Contribution

Fork or issue or pull request.
Everything ok.


## Licence

[The MIT License (MIT) | Open Source Initiative](http://opensource.org/licenses/MIT)


## Author

[peketamin](https://github.com/peketamin)


## Resources

* [Product Advertising API](https://affiliate.amazon.co.jp/gp/advertising/api/detail/main.html) (ja: amazon.co.jp)
* [Product Advertising API](https://affiliate-program.amazon.com/gp/advertising/api/detail/main.html) (en: amazon.com)
