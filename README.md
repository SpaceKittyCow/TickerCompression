Ticker Compression

Ticker Stock Pricing Compression Module 

Used to compress and decopress JSON file from Polygon.io /v2/ticks/stocks/trades/{ticker} endpoint, though with some alteration could be used for any JSON that has some kind of mapping in it. 

Would recommend not running the whole day, but just a constrained results call, like the example below. It take forever for a day long call to run.

Sample call
./polygon.io -a={apiKey} -d=14 -m=10 -y=2020 -r=40
