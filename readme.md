Concurrent-safe string to id mapping. There's no internal state / consistency. If "leto" was  initially assigned 4 and the application restarts, it'll likely get assigned a different value. Furthermore, if "leto" is removed and then re-added, it'll get a new id.

What's the point? The point is to be able to expose string ids ("leto" or guids) but interally use more efficient integers. If the internal integer changes, that's fine because it has no meaning on its own

````go
// The parameter is the # of buckets to use. Buckets help shard write-locks.
// Write-locks are short-lived, so there should be no reason for this to be
// very large.
map := idmap.New(4)


// the 2nd parameter tells Get to create the mapping if it doesn't exist
// id1 will be equal to 1
id1 := map.Get("leto", true)

// id2 wil be equal to 2
id2 := map.Get("ghanima", true)

map.Remove("leto")

// id3 will be equal to 3
id3 := map.Get("paul", true)

// id4 will be equal to 0 (create is false)
id4 := map.Get("jessica", false)
````
