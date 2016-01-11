Request Documentation
=====================

## Go Library ##

For some reason, the goamz/s3 library isn't sending requests out. It isn't connecting to `objects.liquidweb.services` at all, and isn't even trying. I haven't dug into this - perhaps I have to initialize the connection?

## awscli ##

Not like this matters anyway - I installed awscli from brew, and tried using it with:

```
jack@jack-mobile:~|⇒  aws s3 ls --endpoint-url https://objects.liquidweb.services

A client error (InvalidArgument) occurred when calling the ListBuckets operation: Unknown
```

The Request and response sent back and forth was:

```
GET / HTTP/1.1
Host: objects.liquidweb.services
Accept-Encoding: identity
X-Amz-Content-SHA256: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Authorization: AWS4-HMAC-SHA256 Credential=W2JPSSVWUJOBLFE4JYER/20150401/LiquidWeb/s3/aws4_request, SignedHeaders=host;user-agent;x-amz-content-sha256;x-amz-date, Signature=169e27faecbad7f89f75a4b11b4cbda38a1e9d68ded959069b11629ff827c585
X-Amz-Date: 20150401T182532Z
User-Agent: aws-cli/1.7.0 Python/2.7.6 Darwin/14.3.0
```

```
HTTP/1.1 400 Bad Request
Date: Wed, 01 Apr 2015 18:25:33 GMT
Server: Apache/2.2.22 (Ubuntu)
Accept-Ranges: bytes
Content-Length: 81
Content-Type: application/xml
Set-Cookie: RADOSGWLB=; Expires=Thu, 01-Jan-1970 00:00:01 GMT; path=/

<?xml version="1.0" encoding="UTF-8"?><Error><Code>InvalidArgument</Code></Error>
```

That's just to list the buckets - not even doing anything. Similarlly:

```
jack@jack-mobile:~|⇒  aws s3api --endpoint-url https://objects.liquidweb.services list-buckets

A client error (InvalidArgument) occurred when calling the ListBuckets operation: Unknown
```

```
GET / HTTP/1.1
Host: objects.liquidweb.services
Accept-Encoding: identity
X-Amz-Content-SHA256: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
Authorization: AWS4-HMAC-SHA256 Credential=W2JPSSVWUJOBLFE4JYER/20150401/LiquidWeb/s3/aws4_request, SignedHeaders=host;user-agent;x-amz-content-sha256;x-amz-date, Signature=92f9a95402d1fb4b06fdb2693ccef4a98dc09dd873f612829417172896832b3a
X-Amz-Date: 20150401T183305Z
User-Agent: aws-cli/1.7.0 Python/2.7.6 Darwin/14.3.0
```

```
HTTP/1.1 400 Bad Request
Date: Wed, 01 Apr 2015 18:33:05 GMT
Server: Apache/2.2.22 (Ubuntu)
Accept-Ranges: bytes
Content-Length: 81
Content-Type: application/xml
Set-Cookie: RADOSGWLB=; Expires=Thu, 01-Jan-1970 00:00:01 GMT; path=/

<?xml version="1.0" encoding="UTF-8"?><Error><Code>InvalidArgument</Code></Error>
```

## Cyberduck ##

I've found that with CyberDuck, I **can** request objects from our object storage, and send them back. Cyberduck is sending these requests to the url `objects.liquidweb.com` over port `443`.

```
GET / HTTP/1.1
Date: Wed, 01 Apr 2015 18:26:05 GMT
Authorization: AWS W2JPSSVWUJOBLFE4JYER:kswQmIEYB7XbOQLLSxrYWL1qARY=
Host: objects.liquidweb.services:443
Connection: Keep-Alive
User-Agent: Cyberduck/4.5.2 (Mac OS X/10.10.3) (x86_64)
Accept-Encoding: gzip,deflate
```

```
HTTP/1.1 200 OK
Date: Wed, 01 Apr 2015 18:26:06 GMT
Server: Apache/2.2.22 (Ubuntu)
Transfer-Encoding: chunked
Content-Type: application/xml
Set-Cookie: RADOSGWLB=; Expires=Thu, 01-Jan-1970 00:00:01 GMT; path=/
Cache-control: private

<?xml version="1.0" encoding="UTF-8"?><ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Owner><ID>QJBZ6R</ID><DisplayName>QJBZ6R</DisplayName></Owner><Buckets><Bucket><Name>stuff</Name><CreationDate>2015-03-31T06:11:55.000Z</CreationDate></Bucket></Buckets></ListAllMyBucketsResult>
```

That appears to rely on a pre-shared key of some sort though, or a hashed session code? I don't see my secretKey in there anywhere. Either way, once that was sent, listing of contents then also worked.

Recreated the connection, and although the string after the accessKey changed, authentication still isn't happening

## Listing of Bucket Contents ##

Once you're using the bucket, it looks for sure like you move over to the url:

`bucketname.objects.liquidweb.services`

```
GET /?delimiter=%2F&max-keys=1000&prefix HTTP/1.1
Date: Wed, 01 Apr 2015 18:26:08 GMT
Authorization: AWS W2JPSSVWUJOBLFE4JYER:bwVCGeYZRDWECXknPb00ZUaAuG8=
Host: stuff.objects.liquidweb.services:443
Connection: Keep-Alive
User-Agent: Cyberduck/4.5.2 (Mac OS X/10.10.3) (x86_64)
Accept-Encoding: gzip,deflate
```

```
HTTP/1.1 200 OK
Date: Wed, 01 Apr 2015 18:26:08 GMT
Server: Apache/2.2.22 (Ubuntu)
Transfer-Encoding: chunked
Content-Type: application/xml
Set-Cookie: RADOSGWLB=; Expires=Thu, 01-Jan-1970 00:00:01 GMT; path=/
Cache-control: private

<?xml version="1.0" encoding="UTF-8"?><ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Name>stuff</Name><Prefix></Prefix><Marker></Marker><MaxKeys>1000</MaxKeys><Delimiter>/</Delimiter><IsTruncated>false</IsTruncated><Contents><Key>sakila-db.zip</Key><LastModified>2015-03-31T06:13:40.000Z</LastModified><ETag>&quot;4e426c9fae8b10c3c8da7198bfaaddc4&quot;</ETag><Size>718438</Size><StorageClass>STANDARD</StorageClass><Owner><ID>QJBZ6R</ID><DisplayName>QJBZ6R</DisplayName></Owner></Contents><CommonPrefixes><Prefix>testfolder/</Prefix></CommonPrefixes></ListBucketResult>
```

```
GET /?delimiter=%2F&max-keys=1000&prefix=testfolder%2F HTTP/1.1
Date: Wed, 01 Apr 2015 18:26:10 GMT
Authorization: AWS W2JPSSVWUJOBLFE4JYER:QMgGlnznWdeXfOjlY/nFg+RRDZU=
Host: stuff.objects.liquidweb.services:443
Connection: Keep-Alive
User-Agent: Cyberduck/4.5.2 (Mac OS X/10.10.3) (x86_64)
Accept-Encoding: gzip,deflate
```

```
HTTP/1.1 200 OK
Date: Wed, 01 Apr 2015 18:26:10 GMT
Server: Apache/2.2.22 (Ubuntu)
Transfer-Encoding: chunked
Content-Type: application/xml
Set-Cookie: RADOSGWLB=; Expires=Thu, 01-Jan-1970 00:00:01 GMT; path=/
Cache-control: private

<?xml version="1.0" encoding="UTF-8"?><ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/"><Name>stuff</Name><Prefix>testfolder/</Prefix><Marker></Marker><MaxKeys>1000</MaxKeys><Delimiter>/</Delimiter><IsTruncated>false</IsTruncated><Contents><Key>testfolder/</Key><LastModified>2015-03-31T06:15:23.000Z</LastModified><ETag>&quot;d41d8cd98f00b204e9800998ecf8427e&quot;</ETag><Size>0</Size><StorageClass>STANDARD</StorageClass><Owner><ID>QJBZ6R</ID><DisplayName>QJBZ6R</DisplayName></Owner></Contents></ListBucketResult>
```


