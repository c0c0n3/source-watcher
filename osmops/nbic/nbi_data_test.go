package nbic

// expired on Wed Sep 08 2021 18:52:11 GMT+0000
var expiredNbiTokenPayload = `{
	"issued_at": 1631123531.1251214,
	"expires": 1631127131.1251214,
	"_id": "TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2",
	"id": "TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2",
	"project_id": "fada443a-905c-4241-8a33-4dcdbdac55e7",
	"project_name": "admin",
	"username": "admin",
	"user_id": "5c6f2d64-9c23-4718-806a-c74c3fc3c98f",
	"admin": true,
	"roles": [{
		"name": "system_admin",
		"id": "cb545e44-cd2b-4c0b-93aa-7e2cee79afc3"
	}]
}`

// expires on Sat May 17 2053 20:38:51 GMT+0000
var validNbiTokenPayload = `{
	"issued_at": 2631127131.1251214,
	"expires": 2631127131.1251214,
	"_id": "TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2",
	"id": "TuD41hLjDvjlR2cPcAFvWcr6FGvRhIk2",
	"project_id": "fada443a-905c-4241-8a33-4dcdbdac55e7",
	"project_name": "admin",
	"username": "admin",
	"user_id": "5c6f2d64-9c23-4718-806a-c74c3fc3c98f",
	"admin": true,
	"roles": [{
		"name": "system_admin",
		"id": "cb545e44-cd2b-4c0b-93aa-7e2cee79afc3"
	}]
}`

var nsDescriptors = `[
    {
        "_id": "aba58e40-d65f-4f4e-be0a-e248c14d3e03",
        "id": "openldap_ns",
        "designer": "OSM",
        "version": "1.0",
        "name": "openldap_ns",
        "vnfd-id": [
            "openldap_knf"
        ],
        "virtual-link-desc": [
            {
                "id": "mgmtnet",
                "mgmt-network": true
            }
        ],
        "df": [
            {
                "id": "default-df",
                "vnf-profile": [
                    {
                        "id": "openldap",
                        "virtual-link-connectivity": [
                            {
                                "constituent-cpd-id": [
                                    {
                                        "constituent-base-element-id": "openldap",
                                        "constituent-cpd-id": "mgmt-ext"
                                    }
                                ],
                                "virtual-link-profile-id": "mgmtnet"
                            }
                        ],
                        "vnfd-id": "openldap_knf"
                    }
                ]
            }
        ],
        "description": "NS consisting of a single KNF openldap_knf connected to mgmt network",
        "_admin": {
            "userDefinedData": {},
            "created": 1631268635.96618,
            "modified": 1631268637.8627107,
            "projects_read": [
                "fada443a-905c-4241-8a33-4dcdbdac55e7"
            ],
            "projects_write": [
                "fada443a-905c-4241-8a33-4dcdbdac55e7"
            ],
            "onboardingState": "ONBOARDED",
            "operationalState": "ENABLED",
            "usageState": "NOT_IN_USE",
            "storage": {
                "fs": "mongo",
                "path": "/app/storage/",
                "folder": "aba58e40-d65f-4f4e-be0a-e248c14d3e03",
                "pkg-dir": "openldap_ns",
                "descriptor": "openldap_ns/openldap_nsd.yaml",
                "zipfile": "openldap_ns.tar.gz"
            }
        },
        "nsdOnboardingState": "ONBOARDED",
        "nsdOperationalState": "ENABLED",
        "nsdUsageState": "NOT_IN_USE",
        "_links": {
            "self": {
                "href": "/nsd/v1/ns_descriptors/aba58e40-d65f-4f4e-be0a-e248c14d3e03"
            },
            "nsd_content": {
                "href": "/nsd/v1/ns_descriptors/aba58e40-d65f-4f4e-be0a-e248c14d3e03/nsd_content"
            }
        }
    },
	{
        "_id": "ddd20a30-d65f-4f4e-be0a-e248c14d3e03",
        "id": "dummy_ns",
        "designer": "OSM",
        "version": "1.0",
        "name": "dummy_ns",
        "vnfd-id": [
            "dummy_knf"
        ],
        "virtual-link-desc": [
            {
                "id": "mgmtnet",
                "mgmt-network": true
            }
        ],
        "df": [
            {
                "id": "default-df",
                "vnf-profile": [
                    {
                        "id": "dummy",
                        "virtual-link-connectivity": [
                            {
                                "constituent-cpd-id": [
                                    {
                                        "constituent-base-element-id": "dummy",
                                        "constituent-cpd-id": "mgmt-ext"
                                    }
                                ],
                                "virtual-link-profile-id": "mgmtnet"
                            }
                        ],
                        "vnfd-id": "dummy_knf"
                    }
                ]
            }
        ],
        "description": "Made-up NS consisting of a single KNF dummy_knf connected to mgmt network",
        "_admin": {
            "userDefinedData": {},
            "created": 1631268635.96618,
            "modified": 1631268637.8627107,
            "projects_read": [
                "fada443a-905c-4241-8a33-4dcdbdac55e7"
            ],
            "projects_write": [
                "fada443a-905c-4241-8a33-4dcdbdac55e7"
            ],
            "onboardingState": "ONBOARDED",
            "operationalState": "ENABLED",
            "usageState": "NOT_IN_USE",
            "storage": {
                "fs": "mongo",
                "path": "/app/storage/",
                "folder": "ddd20a30-d65f-4f4e-be0a-e248c14d3e03",
                "pkg-dir": "dummy_ns",
                "descriptor": "dummy_ns/openldap_nsd.yaml",
                "zipfile": "openldap_ns.tar.gz"
            }
        },
        "nsdOnboardingState": "ONBOARDED",
        "nsdOperationalState": "ENABLED",
        "nsdUsageState": "NOT_IN_USE",
        "_links": {
            "self": {
                "href": "/nsd/v1/ns_descriptors/ddd20a30-d65f-4f4e-be0a-e248c14d3e03"
            },
            "nsd_content": {
                "href": "/nsd/v1/ns_descriptors/ddd20a30-d65f-4f4e-be0a-e248c14d3e03/nsd_content"
            }
        }
    }
]`

var vimAccounts = `[
    {
        "_id": "4a4425f7-3e72-4d45-a4ec-4241186f3547",
        "name": "mylocation1",
        "vim_type": "dummy",
        "description": null,
        "vim_url": "http://localhost/dummy",
        "vim_user": "u",
        "vim_password": "fNnfmd3KFXvfyVKu3nzItg==",
        "vim_tenant_name": "p",
        "_admin": {
            "created": 1631212983.5388303,
            "modified": 1631212983.5388303,
            "projects_read": [
                "fada443a-905c-4241-8a33-4dcdbdac55e7"
            ],
            "projects_write": [
                "fada443a-905c-4241-8a33-4dcdbdac55e7"
            ],
            "operationalState": "ENABLED",
            "operations": [
                {
                    "lcmOperationType": "create",
                    "operationState": "COMPLETED",
                    "startTime": 1631212983.5930278,
                    "statusEnteredTime": 1631212984.0220273,
                    "operationParams": null
                }
            ],
            "current_operation": null,
            "detailed-status": ""
        },
        "schema_version": "1.11",
        "admin": {
            "current_operation": 0
        }
    }
]`
