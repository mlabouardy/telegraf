package couchbase

import (
	"encoding/json"
	"testing"

	"github.com/influxdata/telegraf/testutil"

	"github.com/couchbase/go-couchbase"
)

func TestGatherServer(t *testing.T) {
	var pool couchbase.Pool
	if err := json.Unmarshal([]byte(poolsDefaultResponse), &pool); err != nil {
		t.Fatal("parse poolsDefaultResponse", err)
	}

	if err := json.Unmarshal([]byte(bucketResponse), &pool.BucketMap); err != nil {
		t.Fatal("parse bucketResponse", err)
	}
	var cb Couchbase
	var acc testutil.Accumulator
	cb.gatherServer("mycluster", &acc, &pool)
	acc.AssertContainsTaggedFields(t, "couchbase_node",
		map[string]interface{}{"memory_free": 23181365248.0, "memory_total": 64424656896.0},
		map[string]string{"cluster": "mycluster", "hostname": "172.16.10.187:8091"})
	acc.AssertContainsTaggedFields(t, "couchbase_node",
		map[string]interface{}{"memory_free": 23665811456.0, "memory_total": 64424656896.0},
		map[string]string{"cluster": "mycluster", "hostname": "172.16.10.65:8091"})
	acc.AssertContainsTaggedFields(t, "couchbase_bucket",
		map[string]interface{}{
			"quota_percent_used": 68.85424936294555,
			"ops_per_sec":        5686.789686789687,
			"disk_fetches":       0.0,
			"item_count":         943239752.0,
			"disk_used":          409178772321.0,
			"data_used":          212179309111.0,
			"mem_used":           202156957464.0,
		},
		map[string]string{"cluster": "mycluster", "bucket": "blastro-df"})
}

func TestSanitizeURI(t *testing.T) {

	var sanitizeTest = []struct {
		input    string
		expected string
	}{
		{"http://user:password@localhost:121", "http://localhost:121"},
		{"user:password@localhost:12/endpoint", "localhost:12/endpoint"},
		{"https://mail@address.com:password@localhost", "https://localhost"},
		{"localhost", "localhost"},
		{"user:password@localhost:2321", "localhost:2321"},
		{"http://user:password@couchbase-0.example.com:8091/endpoint", "http://couchbase-0.example.com:8091/endpoint"},
		{" ", " "},
	}

	for _, test := range sanitizeTest {
		result := regexpURI.ReplaceAllString(test.input, "${1}")

		if result != test.expected {
			t.Errorf("TestSanitizeAddress: input %s, expected %s, actual %s", test.input, test.expected, result)
		}
	}
}

// From `/pools/default` on a real cluster
const poolsDefaultResponse string = `{"storageTotals":{"ram":{"total":450972598272,"quotaTotal":360777252864,"quotaUsed":360777252864,"used":446826622976,"usedByData":255061495696,"quotaUsedPerNode":51539607552,"quotaTotalPerNode":51539607552},"hdd":{"total":1108766539776,"quotaTotal":1108766539776,"used":559135126484,"usedByData":515767865143,"free":498944942902}},"serverGroupsUri":"/pools/default/serverGroups?v=98656394","name":"default","alerts":["Metadata overhead warning. Over  63% of RAM allocated to bucket  \"blastro-df\" on node \"172.16.8.148\" is taken up by keys and metadata.","Metadata overhead warning. Over  65% of RAM allocated to bucket  \"blastro-df\" on node \"172.16.10.65\" is taken up by keys and metadata.","Metadata overhead warning. Over  64% of RAM allocated to bucket  \"blastro-df\" on node \"172.16.13.173\" is taken up by keys and metadata.","Metadata overhead warning. Over  65% of RAM allocated to bucket  \"blastro-df\" on node \"172.16.15.75\" is taken up by keys and metadata.","Metadata overhead warning. Over  65% of RAM allocated to bucket  \"blastro-df\" on node \"172.16.13.105\" is taken up by keys and metadata.","Metadata overhead warning. Over  64% of RAM allocated to bucket  \"blastro-df\" on node \"172.16.8.127\" is taken up by keys and metadata.","Metadata overhead warning. Over  63% of RAM allocated to bucket  \"blastro-df\" on node \"172.16.15.120\" is taken up by keys and metadata.","Metadata overhead warning. Over  66% of RAM allocated to bucket  \"blastro-df\" on node \"172.16.10.187\" is taken up by keys and metadata."],"alertsSilenceURL":"/controller/resetAlerts?token=2814&uuid=2bec87861652b990cf6aa5c7ee58c253","nodes":[{"systemStats":{"cpu_utilization_rate":35.43307086614173,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23181365248},"interestingStats":{"cmd_get":17.98201798201798,"couch_docs_actual_disk_size":68506048063,"couch_docs_data_size":38718796110,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":140158886,"curr_items_tot":279374646,"ep_bg_fetched":0.999000999000999,"get_hits":10.98901098901099,"mem_used":36497390640,"ops":829.1708291708292,"vb_replica_curr_items":139215760},"uptime":"341236","memoryTotal":64424656896,"memoryFree":23181365248,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"couchApiBase":"http://172.16.10.187:8092/","clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.10.187","thisNode":true,"hostname":"172.16.10.187:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}},{"systemStats":{"cpu_utilization_rate":47.38255033557047,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23665811456},"interestingStats":{"cmd_get":172.8271728271728,"couch_docs_actual_disk_size":79360565405,"couch_docs_data_size":38736382876,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":140174377,"curr_items_tot":279383025,"ep_bg_fetched":0.999000999000999,"get_hits":167.8321678321678,"mem_used":36650059656,"ops":1685.314685314685,"vb_replica_curr_items":139208648},"uptime":"341210","memoryTotal":64424656896,"memoryFree":23665811456,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"couchApiBase":"http://172.16.10.65:8092/","clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.10.65","hostname":"172.16.10.65:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}},{"systemStats":{"cpu_utilization_rate":25.5586592178771,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23726600192},"interestingStats":{"cmd_get":63.06306306306306,"couch_docs_actual_disk_size":79345105217,"couch_docs_data_size":38728086130,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":139195268,"curr_items_tot":279349113,"ep_bg_fetched":0,"get_hits":53.05305305305306,"mem_used":36476665576,"ops":1878.878878878879,"vb_replica_curr_items":140153845},"uptime":"341210","memoryTotal":64424656896,"memoryFree":23726600192,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"couchApiBase":"http://172.16.13.105:8092/","clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.13.105","hostname":"172.16.13.105:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}},{"systemStats":{"cpu_utilization_rate":26.45803698435277,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23854841856},"interestingStats":{"cmd_get":51.05105105105105,"couch_docs_actual_disk_size":74465931949,"couch_docs_data_size":38723830730,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":139209869,"curr_items_tot":279380019,"ep_bg_fetched":0,"get_hits":47.04704704704704,"mem_used":36471784896,"ops":1831.831831831832,"vb_replica_curr_items":140170150},"uptime":"340526","memoryTotal":64424656896,"memoryFree":23854841856,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"couchApiBase":"http://172.16.13.173:8092/","clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.13.173","hostname":"172.16.13.173:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}},{"systemStats":{"cpu_utilization_rate":47.31034482758621,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23773573120},"interestingStats":{"cmd_get":77.07707707707708,"couch_docs_actual_disk_size":74743093945,"couch_docs_data_size":38594660087,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":139215932,"curr_items_tot":278427644,"ep_bg_fetched":0,"get_hits":53.05305305305305,"mem_used":36306500344,"ops":1981.981981981982,"vb_replica_curr_items":139211712},"uptime":"340495","memoryTotal":64424656896,"memoryFree":23773573120,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"couchApiBase":"http://172.16.15.120:8092/","clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.15.120","hostname":"172.16.15.120:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}},{"systemStats":{"cpu_utilization_rate":17.60660247592847,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23662190592},"interestingStats":{"cmd_get":146.8531468531468,"couch_docs_actual_disk_size":72932847344,"couch_docs_data_size":38581771457,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":139226879,"curr_items_tot":278436540,"ep_bg_fetched":0,"get_hits":144.8551448551448,"mem_used":36421860496,"ops":1495.504495504495,"vb_replica_curr_items":139209661},"uptime":"337174","memoryTotal":64424656896,"memoryFree":23662190592,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"couchApiBase":"http://172.16.8.127:8092/","clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.8.127","hostname":"172.16.8.127:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}},{"systemStats":{"cpu_utilization_rate":21.68831168831169,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":24049729536},"interestingStats":{"cmd_get":11.98801198801199,"couch_docs_actual_disk_size":66414273220,"couch_docs_data_size":38587642702,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":139193759,"curr_items_tot":278398926,"ep_bg_fetched":0,"get_hits":9.990009990009991,"mem_used":36237234088,"ops":883.1168831168832,"vb_replica_curr_items":139205167},"uptime":"341228","memoryTotal":64424656896,"memoryFree":24049729536,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"couchApiBase":"http://172.16.8.148:8092/","clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.8.148","hostname":"172.16.8.148:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}}],"buckets":{"uri":"/pools/default/buckets?v=74117050&uuid=2bec87861652b990cf6aa5c7ee58c253","terseBucketsBase":"/pools/default/b/","terseStreamingBucketsBase":"/pools/default/bs/"},"remoteClusters":{"uri":"/pools/default/remoteClusters?uuid=2bec87861652b990cf6aa5c7ee58c253","validateURI":"/pools/default/remoteClusters?just_validate=1"},"controllers":{"addNode":{"uri":"/controller/addNode?uuid=2bec87861652b990cf6aa5c7ee58c253"},"rebalance":{"uri":"/controller/rebalance?uuid=2bec87861652b990cf6aa5c7ee58c253"},"failOver":{"uri":"/controller/failOver?uuid=2bec87861652b990cf6aa5c7ee58c253"},"startGracefulFailover":{"uri":"/controller/startGracefulFailover?uuid=2bec87861652b990cf6aa5c7ee58c253"},"reAddNode":{"uri":"/controller/reAddNode?uuid=2bec87861652b990cf6aa5c7ee58c253"},"reFailOver":{"uri":"/controller/reFailOver?uuid=2bec87861652b990cf6aa5c7ee58c253"},"ejectNode":{"uri":"/controller/ejectNode?uuid=2bec87861652b990cf6aa5c7ee58c253"},"setRecoveryType":{"uri":"/controller/setRecoveryType?uuid=2bec87861652b990cf6aa5c7ee58c253"},"setAutoCompaction":{"uri":"/controller/setAutoCompaction?uuid=2bec87861652b990cf6aa5c7ee58c253","validateURI":"/controller/setAutoCompaction?just_validate=1"},"clusterLogsCollection":{"startURI":"/controller/startLogsCollection?uuid=2bec87861652b990cf6aa5c7ee58c253","cancelURI":"/controller/cancelLogsCollection?uuid=2bec87861652b990cf6aa5c7ee58c253"},"replication":{"createURI":"/controller/createReplication?uuid=2bec87861652b990cf6aa5c7ee58c253","validateURI":"/controller/createReplication?just_validate=1"},"setFastWarmup":{"uri":"/controller/setFastWarmup?uuid=2bec87861652b990cf6aa5c7ee58c253","validateURI":"/controller/setFastWarmup?just_validate=1"}},"rebalanceStatus":"none","rebalanceProgressUri":"/pools/default/rebalanceProgress","stopRebalanceUri":"/controller/stopRebalance?uuid=2bec87861652b990cf6aa5c7ee58c253","nodeStatusesUri":"/nodeStatuses","maxBucketCount":10,"autoCompactionSettings":{"parallelDBAndViewCompaction":false,"databaseFragmentationThreshold":{"percentage":50,"size":"undefined"},"viewFragmentationThreshold":{"percentage":50,"size":"undefined"}},"fastWarmupSettings":{"fastWarmupEnabled":true,"minMemoryThreshold":10,"minItemsThreshold":10},"tasks":{"uri":"/pools/default/tasks?v=97479372"},"visualSettingsUri":"/internalSettings/visual?v=7111573","counters":{"rebalance_success":4,"rebalance_start":6,"rebalance_stop":2}}`

// From `/pools/default/buckets/blastro-df` on a real cluster
const bucketResponse string = `{"blastro-df": {"name":"blastro-df","bucketType":"membase","authType":"sasl","saslPassword":"","proxyPort":0,"replicaIndex":false,"uri":"/pools/default/buckets/blastro-df?bucket_uuid=2e6b9dc4c278300ce3a4f27ad540323f","streamingUri":"/pools/default/bucketsStreaming/blastro-df?bucket_uuid=2e6b9dc4c278300ce3a4f27ad540323f","localRandomKeyUri":"/pools/default/buckets/blastro-df/localRandomKey","controllers":{"compactAll":"/pools/default/buckets/blastro-df/controller/compactBucket","compactDB":"/pools/default/buckets/default/controller/compactDatabases","purgeDeletes":"/pools/default/buckets/blastro-df/controller/unsafePurgeBucket","startRecovery":"/pools/default/buckets/blastro-df/controller/startRecovery"},"nodes":[{"couchApiBase":"http://172.16.8.148:8092/blastro-df%2B2e6b9dc4c278300ce3a4f27ad540323f","systemStats":{"cpu_utilization_rate":18.39557399723375,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23791935488},"interestingStats":{"cmd_get":10.98901098901099,"couch_docs_actual_disk_size":79525832077,"couch_docs_data_size":38633186946,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":139229304,"curr_items_tot":278470058,"ep_bg_fetched":0,"get_hits":5.994005994005994,"mem_used":36284362960,"ops":1275.724275724276,"vb_replica_curr_items":139240754},"uptime":"343968","memoryTotal":64424656896,"memoryFree":23791935488,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"replication":1,"clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.8.148","hostname":"172.16.8.148:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}},{"couchApiBase":"http://172.16.8.127:8092/blastro-df%2B2e6b9dc4c278300ce3a4f27ad540323f","systemStats":{"cpu_utilization_rate":21.97183098591549,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23533023232},"interestingStats":{"cmd_get":39.96003996003996,"couch_docs_actual_disk_size":63322357663,"couch_docs_data_size":38603481061,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":139262616,"curr_items_tot":278508069,"ep_bg_fetched":0.999000999000999,"get_hits":30.96903096903097,"mem_used":36475078736,"ops":1370.629370629371,"vb_replica_curr_items":139245453},"uptime":"339914","memoryTotal":64424656896,"memoryFree":23533023232,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"replication":1,"clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.8.127","hostname":"172.16.8.127:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}},{"couchApiBase":"http://172.16.15.120:8092/blastro-df%2B2e6b9dc4c278300ce3a4f27ad540323f","systemStats":{"cpu_utilization_rate":23.38028169014084,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23672963072},"interestingStats":{"cmd_get":88.08808808808809,"couch_docs_actual_disk_size":80260594761,"couch_docs_data_size":38632863189,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":139251563,"curr_items_tot":278498913,"ep_bg_fetched":0,"get_hits":74.07407407407408,"mem_used":36348663000,"ops":1707.707707707708,"vb_replica_curr_items":139247350},"uptime":"343235","memoryTotal":64424656896,"memoryFree":23672963072,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"replication":1,"clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.15.120","hostname":"172.16.15.120:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}},{"couchApiBase":"http://172.16.13.173:8092/blastro-df%2B2e6b9dc4c278300ce3a4f27ad540323f","systemStats":{"cpu_utilization_rate":22.15988779803646,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23818825728},"interestingStats":{"cmd_get":103.1031031031031,"couch_docs_actual_disk_size":68247785524,"couch_docs_data_size":38747583467,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":139245453,"curr_items_tot":279451313,"ep_bg_fetched":1.001001001001001,"get_hits":86.08608608608608,"mem_used":36524715864,"ops":1749.74974974975,"vb_replica_curr_items":140205860},"uptime":"343266","memoryTotal":64424656896,"memoryFree":23818825728,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"replication":1,"clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.13.173","hostname":"172.16.13.173:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}},{"couchApiBase":"http://172.16.13.105:8092/blastro-df%2B2e6b9dc4c278300ce3a4f27ad540323f","systemStats":{"cpu_utilization_rate":21.94444444444444,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23721426944},"interestingStats":{"cmd_get":113.1131131131131,"couch_docs_actual_disk_size":68102832275,"couch_docs_data_size":38747477407,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":139230887,"curr_items_tot":279420530,"ep_bg_fetched":0,"get_hits":106.1061061061061,"mem_used":36524887624,"ops":1799.7997997998,"vb_replica_curr_items":140189643},"uptime":"343950","memoryTotal":64424656896,"memoryFree":23721426944,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"replication":1,"clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.13.105","hostname":"172.16.13.105:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}},{"couchApiBase":"http://172.16.10.65:8092/blastro-df%2B2e6b9dc4c278300ce3a4f27ad540323f","systemStats":{"cpu_utilization_rate":60.62176165803109,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23618203648},"interestingStats":{"cmd_get":30.96903096903097,"couch_docs_actual_disk_size":69052175561,"couch_docs_data_size":38755695030,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":140210194,"curr_items_tot":279454253,"ep_bg_fetched":0,"get_hits":26.97302697302698,"mem_used":36543072472,"ops":1337.662337662338,"vb_replica_curr_items":139244059},"uptime":"343950","memoryTotal":64424656896,"memoryFree":23618203648,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"replication":1,"clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.10.65","hostname":"172.16.10.65:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}},{"couchApiBase":"http://172.16.10.187:8092/blastro-df%2B2e6b9dc4c278300ce3a4f27ad540323f","systemStats":{"cpu_utilization_rate":21.83588317107093,"swap_total":0,"swap_used":0,"mem_total":64424656896,"mem_free":23062269952},"interestingStats":{"cmd_get":33.03303303303304,"couch_docs_actual_disk_size":74422029546,"couch_docs_data_size":38758172837,"couch_views_actual_disk_size":0,"couch_views_data_size":0,"curr_items":140194321,"curr_items_tot":279445526,"ep_bg_fetched":0,"get_hits":21.02102102102102,"mem_used":36527676832,"ops":1088.088088088088,"vb_replica_curr_items":139251205},"uptime":"343971","memoryTotal":64424656896,"memoryFree":23062269952,"mcdMemoryReserved":49152,"mcdMemoryAllocated":49152,"replication":1,"clusterMembership":"active","recoveryType":"none","status":"healthy","otpNode":"ns_1@172.16.10.187","thisNode":true,"hostname":"172.16.10.187:8091","clusterCompatibility":196608,"version":"3.0.1-1444-rel-community","os":"x86_64-unknown-linux-gnu","ports":{"proxy":11211,"direct":11210}}],"stats":{"uri":"/pools/default/buckets/blastro-df/stats","directoryURI":"/pools/default/buckets/blastro-df/statsDirectory","nodeStatsListURI":"/pools/default/buckets/blastro-df/nodes"},"ddocs":{"uri":"/pools/default/buckets/blastro-df/ddocs"},"nodeLocator":"vbucket","fastWarmupSettings":false,"autoCompactionSettings":false,"uuid":"2e6b9dc4c278300ce3a4f27ad540323f","vBucketServerMap":{"hashAlgorithm":"CRC","numReplicas":1,"serverList":["172.16.10.187:11210","172.16.10.65:11210","172.16.13.105:11210","172.16.13.173:11210","172.16.15.120:11210","172.16.8.127:11210","172.16.8.148:11210"],"vBucketMap":[[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,6],[0,6],[0,6],[0,6],[0,6],[1,3],[1,3],[1,3],[1,4],[1,4],[1,5],[1,5],[1,5],[1,5],[1,5],[1,5],[1,5],[1,5],[1,5],[1,5],[0,2],[0,2],[0,2],[0,2],[0,2],[0,2],[0,2],[0,2],[0,2],[0,2],[0,2],[0,2],[0,2],[0,2],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[2,3],[2,3],[2,5],[2,5],[2,5],[2,5],[2,5],[2,5],[2,5],[2,5],[2,5],[0,3],[0,3],[0,3],[0,3],[0,3],[0,3],[0,3],[0,3],[0,3],[0,3],[0,3],[0,3],[0,3],[0,3],[2,5],[2,5],[2,6],[2,6],[2,6],[2,6],[2,6],[2,6],[2,6],[2,6],[3,5],[3,5],[3,5],[3,5],[3,5],[3,5],[3,5],[3,5],[3,5],[3,5],[0,4],[0,4],[0,4],[0,4],[0,4],[0,4],[0,4],[0,4],[0,4],[0,4],[0,4],[0,4],[0,4],[0,4],[3,5],[3,5],[3,5],[3,5],[3,5],[3,5],[3,6],[3,6],[3,6],[3,6],[4,5],[4,5],[4,5],[4,5],[4,5],[4,5],[4,5],[4,5],[4,5],[4,5],[0,6],[0,6],[0,6],[0,6],[0,6],[0,6],[0,6],[0,6],[0,6],[0,6],[0,6],[0,6],[0,6],[0,6],[0,6],[5,3],[5,4],[5,4],[5,4],[5,4],[5,4],[5,4],[5,4],[5,4],[5,4],[6,5],[6,5],[6,5],[6,5],[6,5],[6,5],[6,5],[6,5],[6,5],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[0,3],[0,3],[0,3],[0,4],[0,4],[0,4],[0,4],[0,4],[0,4],[0,4],[0,5],[0,5],[0,5],[0,5],[0,5],[0,5],[0,5],[0,5],[0,5],[0,5],[1,2],[1,2],[1,2],[1,2],[1,2],[1,2],[1,2],[1,2],[1,2],[1,2],[1,2],[1,2],[1,2],[1,2],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,5],[2,5],[2,5],[2,5],[2,5],[2,5],[4,5],[4,5],[4,5],[4,5],[1,3],[1,3],[1,3],[1,3],[1,3],[1,3],[1,3],[1,3],[1,3],[1,3],[1,3],[1,3],[1,3],[1,3],[2,6],[2,6],[3,2],[3,2],[3,4],[3,4],[3,4],[3,4],[3,4],[3,4],[3,4],[3,5],[3,5],[3,5],[3,5],[2,0],[2,0],[2,0],[2,0],[2,0],[1,4],[1,4],[1,4],[1,4],[1,4],[1,4],[1,4],[1,4],[1,4],[1,4],[1,4],[1,4],[1,4],[1,4],[3,6],[3,6],[3,6],[3,6],[3,6],[3,6],[4,2],[4,3],[4,3],[4,3],[4,5],[4,5],[4,5],[4,5],[3,0],[3,0],[3,0],[3,0],[3,0],[3,0],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[1,6],[5,4],[5,4],[5,6],[5,6],[5,6],[5,6],[5,6],[5,6],[5,6],[5,6],[6,5],[6,5],[6,5],[6,5],[6,5],[4,0],[4,0],[4,0],[4,0],[4,0],[2,0],[2,0],[2,0],[2,0],[2,0],[2,0],[2,0],[2,0],[2,0],[2,0],[2,0],[2,0],[2,0],[2,0],[0,4],[0,4],[0,4],[0,5],[0,5],[0,5],[0,5],[0,5],[0,5],[0,5],[0,5],[0,5],[0,5],[0,5],[4,5],[4,5],[4,5],[4,5],[4,5],[4,6],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[1,4],[1,4],[1,4],[1,4],[1,4],[1,4],[1,4],[1,4],[1,5],[1,5],[1,5],[1,5],[1,5],[1,5],[1,5],[4,6],[4,6],[4,6],[4,6],[4,6],[2,3],[2,3],[2,3],[2,3],[2,3],[2,3],[2,3],[2,3],[2,3],[2,3],[2,3],[2,3],[2,3],[2,3],[2,3],[3,4],[3,4],[3,4],[3,5],[3,5],[3,5],[3,5],[5,0],[5,0],[5,0],[2,0],[2,0],[3,0],[3,0],[3,0],[5,3],[5,3],[5,3],[5,3],[5,3],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[2,4],[4,3],[4,3],[4,3],[4,3],[4,3],[4,3],[4,3],[4,5],[4,5],[1,0],[3,0],[3,1],[3,1],[3,1],[3,1],[5,4],[5,4],[5,4],[5,4],[5,4],[2,6],[2,6],[2,6],[2,6],[2,6],[2,6],[2,6],[2,6],[2,6],[2,6],[2,6],[2,6],[2,6],[2,6],[5,6],[5,6],[5,6],[6,2],[6,2],[6,3],[6,3],[6,3],[4,0],[4,0],[4,0],[4,0],[4,0],[4,1],[4,1],[4,1],[5,6],[5,6],[5,6],[5,6],[3,0],[3,0],[3,0],[3,0],[3,0],[3,0],[3,0],[3,0],[3,0],[3,0],[3,0],[3,0],[3,0],[3,0],[0,5],[0,5],[0,5],[0,6],[0,6],[0,6],[0,6],[0,6],[0,1],[0,1],[4,6],[4,6],[4,6],[4,6],[5,0],[5,0],[5,0],[5,0],[5,0],[5,0],[3,1],[3,1],[3,1],[3,1],[3,1],[3,1],[3,1],[3,1],[3,1],[3,1],[3,1],[3,1],[3,1],[3,1],[1,5],[1,5],[1,5],[1,5],[1,5],[1,5],[1,5],[1,6],[2,0],[2,0],[5,2],[5,3],[5,3],[5,3],[5,3],[5,1],[5,1],[5,1],[5,1],[5,1],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[2,5],[2,5],[2,5],[2,5],[2,5],[2,5],[2,5],[4,1],[4,1],[4,1],[5,3],[5,3],[5,3],[5,3],[5,3],[2,0],[5,2],[5,2],[5,2],[5,2],[3,4],[3,4],[3,4],[3,4],[3,4],[3,4],[3,4],[3,4],[3,4],[3,4],[3,4],[3,4],[3,4],[3,4],[3,4],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,0],[1,2],[5,4],[5,4],[5,4],[5,4],[5,4],[5,4],[5,4],[5,4],[5,4],[3,6],[3,6],[3,6],[3,6],[3,6],[3,6],[3,6],[3,6],[3,6],[3,6],[3,6],[3,6],[3,6],[3,6],[4,1],[4,1],[5,0],[5,0],[5,0],[5,0],[5,0],[5,0],[5,0],[5,1],[5,6],[5,6],[5,6],[5,6],[5,6],[5,6],[5,6],[5,6],[5,6],[5,6],[4,0],[4,0],[4,0],[4,0],[4,0],[4,0],[4,0],[4,0],[4,0],[4,0],[4,0],[4,0],[4,0],[4,0],[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,1],[0,2],[0,2],[5,0],[5,0],[5,0],[5,0],[5,0],[5,0],[5,0],[5,0],[0,2],[0,2],[4,1],[4,1],[4,1],[4,1],[4,1],[4,1],[4,1],[4,1],[4,1],[4,1],[4,1],[4,1],[4,1],[4,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[2,1],[5,1],[5,1],[5,1],[5,1],[5,1],[5,1],[5,1],[5,1],[5,1],[3,1],[4,2],[4,2],[4,2],[4,2],[4,2],[4,2],[4,2],[4,2],[4,2],[4,2],[4,2],[4,2],[4,2],[4,2],[4,1],[4,1],[4,2],[4,2],[4,2],[6,3],[6,3],[6,3],[6,3],[6,3],[5,2],[5,2],[5,2],[5,2],[5,2],[5,2],[5,2],[5,2],[5,2],[5,2],[4,3],[4,3],[4,3],[4,3],[4,3],[4,3],[4,3],[4,3],[4,3],[4,3],[4,3],[4,3],[4,3],[4,3],[4,3],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[5,3],[5,3],[5,3],[5,3],[5,3],[5,3],[5,3],[5,3],[5,3],[4,6],[4,6],[4,6],[4,6],[4,6],[4,6],[4,6],[4,6],[4,6],[4,6],[4,6],[4,6],[4,6],[4,6],[5,1],[5,1],[5,1],[5,1],[5,1],[5,1],[5,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,2],[6,2],[6,2],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[6,0],[1,2],[1,2],[1,2],[1,2],[1,2],[1,2],[1,2],[2,1],[2,3],[2,3],[1,2],[1,2],[1,2],[1,3],[1,3],[1,3],[1,3],[1,3],[3,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[6,1],[3,1],[3,1],[3,1],[3,1],[4,2],[4,2],[4,2],[4,2],[4,2],[4,2],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[3,2],[6,3],[6,3],[6,2],[6,2],[6,2],[6,2],[6,2],[6,2],[6,2],[6,2],[6,2],[6,2],[6,2],[6,2],[6,2],[6,2],[5,1],[5,1],[5,2],[5,2],[5,2],[5,2],[5,2],[5,2],[5,2],[5,2],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[5,2],[6,2],[6,2],[6,3],[6,3],[6,3],[6,3],[6,3],[6,3],[6,3],[6,3],[6,3],[6,3],[6,3],[6,3],[6,3],[6,3],[0,2],[0,2],[0,2],[0,2],[0,2],[0,2],[0,2],[0,3],[1,3],[1,3],[6,2],[6,2],[0,3],[0,3],[0,3],[0,3],[0,3],[0,3],[1,3],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,4],[6,5],[6,5],[2,3],[2,3],[2,3],[2,3],[2,3],[2,3],[6,5],[6,5],[6,5],[6,5],[6,5],[6,5],[6,5],[6,5],[6,5],[6,2]]},"replicaNumber":1,"threadsNumber":3,"quota":{"ram":293601280000,"rawRAM":41943040000},"basicStats":{"quotaPercentUsed":68.85424936294555,"opsPerSec":5686.789686789687,"diskFetches":0,"itemCount":943239752,"diskUsed":409178772321,"dataUsed":212179309111,"memUsed":202156957464},"evictionPolicy":"valueOnly","bucketCapabilitiesVer":"","bucketCapabilities":["cbhello","touch","couchapi","cccp","xdcrCheckpointing","nodesExt"]}}`
