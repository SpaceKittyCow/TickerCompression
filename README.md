Ticker Compression

Ticker Stock Pricing Compression Module 

Used to compress and decopress JSON file from Polygon.io /v2/ticks/stocks/trades/{ticker} endpoint, though with some alteration could be used for any JSON that has some kind of mapping in it. 

Would recommend not running the whole day, but just a constrained results call, like the example below. It take forever for a day long call to run.

Logrithimitically saves the data, so the larger the set, the smaller the save file, starting at 40% compression for 10 results onwards.

Sample call:
./polygon.io -a={apiKey} -d=14 -m=10 -y=2020 -r=40

polygon.io --help: 

Usage of ./polygon.io:
  -a string
    	Polygon.io API Key
  -c string
    	Specify a save location of the compressed file (default "./compressedfile")
  -d int
    	Specify a day of the month: 1, 2, 3 ...
  -j string
    	Specify a load location of an orignal file
  -m int
    	Specify a month number: 1 for January, 2 for Feburary ...
  -r int
    	Will go through one pass through instead of all day. Reccomended due to getting the whole day taking forever. Upto 50000
  -s string
    	Specify a save location of the orignal file (default "./orignalfile.json")
  -y int
    	Specify a year: 2021, 2020 ...

