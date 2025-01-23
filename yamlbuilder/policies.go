package yamlbuilder

// the ValidatingAdmissionPolicies to be validated and to be generated into VAP yaml after validation.


var ValidatingAdmissionPolicies  = []FuncInput{
	{
		Name: "etcd-validation-deny",
		Validations: []ValidationInfo{
			{
				Expression: "has(object.spec.etcd.resources.requests)",
				Message:  "The Requests field in etcd.spec.etcd.resources cannot be empty. ",
			},
			{
				Expression: "!has(object.spec.etcd.etcdDefragTimeout)|| string(object.spec.etcd.etcdDefragTimeout).matches('^(([0-9]+(.[0-9]+)?h)?(([0-9]+(.[0-9]+)?m)?(([0-9]+(.[0-9]+)?s)?))?)$')",
				Message: "Invalid Duration given for etcd.spec.etcd.etcdDefragTimeout",
			} ,
			{
				Expression: "!has(object.spec.etcd.heartbeatDuration)|| string(object.spec.etcd.heartbeatDuration).matches('^(([0-9]+(.[0-9]+)?h)?(([0-9]+(.[0-9]+)?m)?(([0-9]+(.[0-9]+)?s)?))?)$')",
				Message: "Invalid Duration given for etcd.spec.etcd.heartbeatDuration",
			} ,
			{
				Expression: "!has(object.spec.etcd.defragmentationSchedule) || object.spec.etcd.defragmentationSchedule.matches('^(\\\\*|([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])|\\\\*/([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])) (\\\\*|([0-9]|1[0-9]|2[0-3])|\\\\*/([0-9]|1[0-9]|2[0-4])) (\\\\*|([1-9]|1[0-9]|2[0-9]|3[0-1])|\\\\*/([1-9]|1[0-9]|2[0-9]|3[0-1])) (\\\\*|([1-9]|1[0-2])|\\\\*/([1-9]|1[0-2])) (\\\\*|([0-6])|\\\\*/([0-6]))$')",
				Message: "Invalid cron expression given for defragmentation schedule (etcd.spec.etcd.defragmentationSchedule)",
			} ,
			{
				Expression: "!has(object.spec.backup.fullSnapshotSchedule) || object.spec.backup.fullSnapshotSchedule.matches('^(\\\\*|([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])|\\\\*/([0-9]|1[0-9]|2[0-9]|3[0-9]|4[0-9]|5[0-9])) (\\\\*|([0-9]|1[0-9]|2[0-3])|\\\\*/([0-9]|1[0-9]|2[0-4])) (\\\\*|([1-9]|1[0-9]|2[0-9]|3[0-1])|\\\\*/([1-9]|1[0-9]|2[0-9]|3[0-1])) (\\\\*|([1-9]|1[0-2])|\\\\*/([1-9]|1[0-2])) (\\\\*|([0-6])|\\\\*/([0-6]))$')",
				Message: "Invalid cron expression given for fullSnapshotSchedule (etcd.spec.backup.fullSnapshotSchedule).",
			} ,
			{
				Expression: "!has(object.spec.backup.garbageCollectionPeriod) || string(object.spec.backup.garbageCollectionPeriod).matches('^(([0-9]+(.[0-9]+)?h)?(([0-9]+(.[0-9]+)?m)?(([0-9]+(.[0-9]+)?s)?))?)$')",
				Message: "Invalid duration given for etcd.spec.backup.garbageCollectionPeriod",
			} ,
			{
				Expression: "!has(object.spec.backup.deltaSnapshotPeriod) || string(object.spec.backup.deltaSnapshotPeriod).matches('^(([0-9]+(.[0-9]+)?h)?(([0-9]+(.[0-9]+)?m)?(([0-9]+(.[0-9]+)?s)?))?)$')",
				Message: "Invalid duration given for etcd.spec.backup.deltaSnapshotPeriod",
			} ,
			{
				Expression: "!(has(object.spec.backup.deltaSnapshotPeriod) && has(object.spec.backup.garbageCollectionPeriod)) || duration(object.spec.backup.deltaSnapshotPeriod).getSeconds() < duration(object.spec.backup.garbageCollectionPeriod).getSeconds()",
				Message: "GarbageCollectionPeriod must be greater than the deltasnapshotperiod",
			} ,
			{
				Expression: "!has(object.spec.backup.etcdSnapshotTimeout) || string(object.spec.backup.etcdSnapshotTimeout).matches('^(([0-9]+(.[0-9]+)?h)?(([0-9]+(.[0-9]+)?m)?(([0-9]+(.[0-9]+)?s)?))?)$')",
				Message: "Invalid duration given for etcd.spec.backup.etcdSnapshotTimeout",
			} ,
			{
				Expression: "!has(object.spec.backup.leaderElection.reelectionPeriod) || string(object.spec.backup.leaderElection.reelectionPeriod).matches('^(([0-9]+(.[0-9]+)?h)?(([0-9]+(.[0-9]+)?m)?(([0-9]+(.[0-9]+)?s)?))?)$')",
				Message: "Invalid duration given for etcd.spec.backup.leaderElection.reelectionPeriod",
			} ,
			{
				Expression: "!has(object.spec.backup.leaderElection.etcdConnectionTimeout) || string(object.spec.backup.leaderElection.etcdConnectionTimeout).matches('^(([0-9]+(.[0-9]+)?h)?(([0-9]+(.[0-9]+)?m)?(([0-9]+(.[0-9]+)?s)?))?)$')",
				Message: "Invalid duration given for etcd.spec.backup.leaderElection.etcdConnectionTimeout",
			},
		},
	},
	{
		Name: "etcd-validation-warn",
		Validations: []ValidationInfo{
			{
				Expression: "object.spec.replicas % 2 == 1",
				Message: "replicas should be odd in number",
			},
			{
				Expression: "!has(object.spec.etcd.quota) || quantity(object.spec.etcd.quota).compareTo(quantity('8Gi')) <= 0",
				Message: "Max size is 8Gi",
			},
		},
	},
}


