I was recovered my photos from formatted drive. However, many of the photos are duplicate, already exist on my master folder (of photos). 

I want to find these duplicates and delete it, and use several photo management software to find the duplicate but none useful. The problem was, they usually find the exact file size to the files. However, my recovered photos has different size because there are several zero bytes concatenated in the tail of the files. 

Thus I have to build my one duplicate finder, by hashing only the front of the photos and then compare to my photos master.