# `zc2` cluster
custom (vCPU: 0, Memory: 712 GB) * 28

**System reserved:** 36 GB (712 GB * 0.05 or 4 GB, whichever is bigger)

**Allocated to the Docker container:** 712 GB - 36 GB = 676 GB [[docker-stack-java.yaml](docker-stack-java.yaml)] [[docker-stack-native.yaml](docker-stack-native.yaml)]

**For Java clusters:**
* `HeapSize` = `ContainerMemory` (676 GB) * 0.9 = 608 GB (`-Xmx` and `-Xms` in [[coordinator jvm.config](coordinator/jvm.config)] and [[worker jvm.config](workers/jvm.config)])
* Presto: [[coordinator config.properties](coordinator/config.properties)] and [[worker config.properties](worker/config.properties)]
  * `query.max-total-memory-per-node` = `HeapSize` (608 GB) * 0.8 = 486 GB
  * `query.max-memory-per-node` = `query.max-total-memory-per-node` (486 GB) * 0.9 = 437 GB
  * `query.max-total-memory` = `query.max-total-memory-per-node` (486 GB) * `[number of nodes]` (28) = 13608 GB
  * `query.max-memory` = `query.max-memory-per-node` (437 GB) * `[number of nodes]` (28) = 12236 GB [[documentation](https://prestodb.io/docs/current/admin/properties.html#memory-management-properties)]
  * `memory.heap-headroom-per-node` = `HeapSize` (608 GB) * 0.2 = 122 GB

**For Prestissimo clusters:**
* Coordinator heap setting same as Java cluster
* `system-memory-gb` = `ContainerMemory` (676 GB) * 0.9 = 608 GB
* `query.max-memory-per-node` = `system-memory-gb` (608 GB) * 1 = 608 GB
* `query-memory-gb` = `query.max-memory-per-node` = 608 GB
  * `MemoryForSpillingAndCaching` = `query.max-memory-per-node` - `query-memory-gb`. We don't need this for the benchmarking now so `query-memory-gb` = `query.max-memory-per-node`