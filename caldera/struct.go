package caldera

import "time"

// OperationReport.Steps {[paw]:[]Steps, [paw]:[]Steps}
type OperationReport struct {
	Name       string        `json:"name"`
	Start      string        `json:"start"`
	Host_group []interface{} `json:"host_group"`
	Steps      interface{}   `json:"steps"`
	Finish     bool          `json:"finish"`
	Planner    string
	Adversary  interface{}
	Jitter     string
	Objectives interface{}
	Facts      []interface{}
}

type Step struct {
	LinkID           string    `json:"link_id"`
	AbilityID        string    `json:"ability_id"`
	Command          string    `json:"command"`
	PlaintextCommand string    `json:"plaintext_command"`
	Delegated        time.Time `json:"delegated"`
	Run              time.Time `json:"run"`
	Status           int       `json:"status"`
	Platform         string    `json:"platform"`
	Executor         string    `json:"executor"`
	Pid              int       `json:"pid"`
	Description      string    `json:"description"`
	Name             string    `json:"name"`
	Attack           struct {
		Tactic        string `json:"tactic"`
		TechniqueName string `json:"technique_name"`
		TechniqueID   string `json:"technique_id"`
	} `json:"attack"`
	Output string `json:"output"`
}

// {
// 	"objective": {
// 		"name": "default",
// 		"goals": [
// 			{
// 				"operator": "==",
// 				"count": 1048576,
// 				"value": "complete",
// 				"achieved": false,
// 				"target": "exhaustion"
// 			}
// 		],
// 		"description": "This is a default objective that runs forever.",
// 		"id": "495a9828-cab1-44dd-a0ca-66e58177d8cc",
// 		"percentage": 0.0
// 	},
// 	"chain": [
// 		{
// 			"decide": "2022-12-08T06:16:34Z",
// 			"deadman": false,
// 			"score": 0,
// 			"visibility": {
// 				"score": 50,
// 				"adjustments": []
// 			},
// 			"plaintext_command": "dGFza2xpc3Q=",
// 			"paw": "ydhydh11111",
// 			"used": [],
// 			"executor": {
// 				"name": "cmd",
// 				"platform": "windows",
// 				"parsers": [],
// 				"variations": [],
// 				"uploads": [],
// 				"build_target": null,
// 				"code": null,
// 				"payloads": [],
// 				"additional_info": {},
// 				"command": "tasklist",
// 				"language": null,
// 				"cleanup": [],
// 				"timeout": 60
// 			},
// 			"host": "DESKTOP-48JI9HA",
// 			"cleanup": 0,
// 			"facts": [],
// 			"unique": "59781280-ac63-4e72-8eae-95a92cc0a9d0",
// 			"pid": "",
// 			"relationships": [],
// 			"status": -3,
// 			"collect": "",
// 			"pin": 0,
// 			"id": "59781280-ac63-4e72-8eae-95a92cc0a9d0",
// 			"finish": "",
// 			"command": "dGFza2xpc3Q=",
// 			"output": "False",
// 			"ability": {
// 				"technique_name": "auto-generated",
// 				"privilege": null,
// 				"requirements": [],
// 				"ability_id": "e2f8bf00-09d7-48cc-b673-f912ea016d69",
// 				"name": "tasklist",
// 				"repeatable": null,
// 				"technique_id": "auto-generated",
// 				"singleton": null,
// 				"executors": [
// 					{
// 						"name": "cmd",
// 						"platform": "windows",
// 						"parsers": [],
// 						"variations": [],
// 						"uploads": [],
// 						"build_target": null,
// 						"code": null,
// 						"payloads": [],
// 						"additional_info": {},
// 						"command": "tasklist",
// 						"language": null,
// 						"cleanup": [],
// 						"timeout": 60
// 					}
// 				],
// 				"description": "Manual command ability",
// 				"additional_info": {},
// 				"access": {},
// 				"plugin": null,
// 				"delete_payload": null,
// 				"buckets": [],
// 				"tactic": "auto-generated"
// 			}
// 		}
// 	],
// 	"use_learning_parsers": true,
// 	"host_group": [
// 		{
// 			"architecture": "amd64",
// 			"upstream_dest": "http://pdxf.tk:8888",
// 			"proxy_receivers": {},
// 			"watchdog": 0,
// 			"host_ip_addrs": [
// 				"192.168.122.244"
// 			],
// 			"group": "ydhydh",
// 			"trusted": false,
// 			"executors": [
// 				"cmd",
// 				"psh",
// 				"proc"
// 			],
// 			"deadman_enabled": true,
// 			"ppid": 3408,
// 			"sleep_min": 30,
// 			"location": "C:\\Users\\Public\\splunkd.exe",
// 			"pending_contact": "HTTP",
// 			"server": "http://pdxf.tk:8888",
// 			"links": [
// 				{
// 					"visibility": {
// 						"score": 50,
// 						"adjustments": []
// 					},
// 					"plaintext_command": "Q2xlYXItSGlzdG9yeTtDbGVhcg==",
// 					"used": [],
// 					"executor": {
// 						"name": "psh",
// 						"platform": "windows",
// 						"parsers": [],
// 						"variations": [],
// 						"uploads": [],
// 						"build_target": null,
// 						"code": null,
// 						"payloads": [],
// 						"additional_info": {},
// 						"command": "Clear-History;Clear",
// 						"language": null,
// 						"cleanup": [],
// 						"timeout": 60
// 					},
// 					"facts": [],
// 					"collect": "2022-12-07T06:13:01Z",
// 					"pin": 0,
// 					"jitter": 0,
// 					"command": "Q2xlYXItSGlzdG9yeTtDbGVhcg==",
// 					"output": "False",
// 					"ability": {
// 						"technique_name": "Indicator Removal on Host: Clear Command History",
// 						"privilege": null,
// 						"requirements": [],
// 						"ability_id": "43b3754c-def4-4699-a673-1d85648fda6a",
// 						"name": "Avoid logs",
// 						"repeatable": false,
// 						"technique_id": "T1070.003",
// 						"singleton": false,
// 						"executors": [
// 							{
// 								"name": "sh",
// 								"platform": "darwin",
// 								"parsers": [],
// 								"variations": [],
// 								"uploads": [],
// 								"build_target": null,
// 								"code": null,
// 								"payloads": [],
// 								"additional_info": {},
// 								"command": "> $HOME/.bash_history && unset HISTFILE",
// 								"language": null,
// 								"cleanup": [],
// 								"timeout": 60
// 							},
// 							{
// 								"name": "sh",
// 								"platform": "linux",
// 								"parsers": [],
// 								"variations": [],
// 								"uploads": [],
// 								"build_target": null,
// 								"code": null,
// 								"payloads": [],
// 								"additional_info": {},
// 								"command": "> $HOME/.bash_history && unset HISTFILE",
// 								"language": null,
// 								"cleanup": [],
// 								"timeout": 60
// 							},
// 							{
// 								"name": "psh",
// 								"platform": "windows",
// 								"parsers": [],
// 								"variations": [],
// 								"uploads": [],
// 								"build_target": null,
// 								"code": null,
// 								"payloads": [],
// 								"additional_info": {},
// 								"command": "Clear-History;Clear",
// 								"language": null,
// 								"cleanup": [],
// 								"timeout": 60
// 							}
// 						],
// 						"description": "Stop terminal from logging history",
// 						"additional_info": {},
// 						"access": {},
// 						"plugin": "stockpile",
// 						"delete_payload": true,
// 						"buckets": [
// 							"defense-evasion"
// 						],
// 						"tactic": "defense-evasion"
// 					},
// 					"decide": "2022-12-07T06:13:01Z",
// 					"paw": "ydhydh11111",
// 					"deadman": false,
// 					"host": "DESKTOP-48JI9HA",
// 					"agent_reported_time": "2022-12-07T06:09:26Z",
// 					"cleanup": 0,
// 					"unique": "fc659185-4676-4539-9af3-a9e1a3120f7d",
// 					"pid": "6076",
// 					"relationships": [],
// 					"status": 0,
// 					"id": "fc659185-4676-4539-9af3-a9e1a3120f7d",
// 					"finish": "2022-12-07T06:13:01Z",
// 					"score": 0
// 				}
// 			],
// 			"username": "DESKTOP-48JI9HA\\pdxf",
// 			"proxy_chain": [],
// 			"origin_link_id": "",
// 			"paw": "ydhydh11111",
// 			"privilege": "User",
// 			"host": "DESKTOP-48JI9HA",
// 			"last_seen": "2022-12-07T06:26:42Z",
// 			"pid": 2144,
// 			"sleep_max": 60,
// 			"platform": "windows",
// 			"display_name": "DESKTOP-48JI9HA$DESKTOP-48JI9HA\\pdxf",
// 			"contact": "HTTP",
// 			"available_contacts": [
// 				"HTTP"
// 			],
// 			"exe_name": "splunkd.exe",
// 			"created": "2022-12-07T06:13:01Z"
// 		}
// 	],
// 	"name": "runtest",
// 	"planner": {
// 		"params": {},
// 		"name": "atomic",
// 		"allow_repeatable_abilities": false,
// 		"stopping_conditions": [],
// 		"module": "plugins.stockpile.app.atomic",
// 		"id": "aaa7c857-37a0-4c4a-85f7-4e9f7f30e31a",
// 		"ignore_enforcement_modules": [],
// 		"description": "During each phase of the operation, the atomic planner iterates through each agent and sends the next\navailable ability it thinks that agent can complete. This decision is based on the agent matching the operating\nsystem (execution platform) of the ability and the ability command having no unsatisfied variables.\nThe planner then waits for each agent to complete its command before determining the subsequent abilities.\nThe abilities are processed in the order set by each agent's atomic ordering.\nFor instance, if agent A has atomic ordering (A1, A2, A3) and agent B has atomic ordering (B1, B2, B3), then\nthe planner would send (A1, B1) in the first phase, then (A2, B2), etc.\n",
// 		"plugin": "stockpile"
// 	},
// 	"source": {
// 		"name": "basic",
// 		"relationships": [],
// 		"rules": [
// 			{
// 				"trait": "file.sensitive.extension",
// 				"match": ".*",
// 				"action": "DENY"
// 			},
// 			{
// 				"trait": "file.sensitive.extension",
// 				"match": "png",
// 				"action": "ALLOW"
// 			},
// 			{
// 				"trait": "file.sensitive.extension",
// 				"match": "yml",
// 				"action": "ALLOW"
// 			},
// 			{
// 				"trait": "file.sensitive.extension",
// 				"match": "wav",
// 				"action": "ALLOW"
// 			}
// 		],
// 		"id": "ed32b9c3-9593-4c33-b0db-e2007315096b",
// 		"plugin": "stockpile",
// 		"adjustments": [],
// 		"facts": [
// 			{
// 				"name": "file.sensitive.extension",
// 				"origin_type": "SEEDED",
// 				"relationships": [],
// 				"technique_id": null,
// 				"source": "ed32b9c3-9593-4c33-b0db-e2007315096b",
// 				"trait": "file.sensitive.extension",
// 				"limit_count": -1,
// 				"value": "wav",
// 				"collected_by": [],
// 				"links": [],
// 				"score": 1,
// 				"created": "2022-12-08T06:13:18Z",
// 				"unique": "file.sensitive.extensionwav"
// 			},
// 			{
// 				"name": "file.sensitive.extension",
// 				"origin_type": "SEEDED",
// 				"relationships": [],
// 				"technique_id": null,
// 				"source": "ed32b9c3-9593-4c33-b0db-e2007315096b",
// 				"trait": "file.sensitive.extension",
// 				"limit_count": -1,
// 				"value": "yml",
// 				"collected_by": [],
// 				"links": [],
// 				"score": 1,
// 				"created": "2022-12-08T06:13:18Z",
// 				"unique": "file.sensitive.extensionyml"
// 			},
// 			{
// 				"name": "file.sensitive.extension",
// 				"origin_type": "SEEDED",
// 				"relationships": [],
// 				"technique_id": null,
// 				"source": "ed32b9c3-9593-4c33-b0db-e2007315096b",
// 				"trait": "file.sensitive.extension",
// 				"limit_count": -1,
// 				"value": "png",
// 				"collected_by": [],
// 				"links": [],
// 				"score": 1,
// 				"created": "2022-12-08T06:13:18Z",
// 				"unique": "file.sensitive.extensionpng"
// 			},
// 			{
// 				"name": "server.malicious.url",
// 				"origin_type": "SEEDED",
// 				"relationships": [],
// 				"technique_id": null,
// 				"source": "ed32b9c3-9593-4c33-b0db-e2007315096b",
// 				"trait": "server.malicious.url",
// 				"limit_count": -1,
// 				"value": "keyloggedsite.com",
// 				"collected_by": [],
// 				"links": [],
// 				"score": 1,
// 				"created": "2022-12-08T06:13:18Z",
// 				"unique": "server.malicious.urlkeyloggedsite.com"
// 			}
// 		]
// 	},
// 	"visibility": 50,
// 	"adversary": {
// 		"name": "ad-hoc",
// 		"tags": [],
// 		"has_repeatable_abilities": false,
// 		"adversary_id": "ad-hoc",
// 		"description": "an empty adversary profile",
// 		"plugin": null,
// 		"objective": "495a9828-cab1-44dd-a0ca-66e58177d8cc",
// 		"atomic_ordering": []
// 	},
// 	"start": "2022-12-08T06:13:18Z",
// 	"auto_close": false,
// 	"jitter": "2/8",
// 	"autonomous": 1,
// 	"id": "829a475b-a486-45f7-8f13-ca84c5b51dff",
// 	"state": "finished",
// 	"group": "ydhydh",
// 	"obfuscator": "plain-text"
// }
