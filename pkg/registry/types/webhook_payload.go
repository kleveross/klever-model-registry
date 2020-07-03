/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package types

// !!! important
// this file content if copy from https://github.com/goharbor/harbor/blob/master/src/pkg/notifier/model/event.go
// if harbor update it, we must update synchronize

// Payload of notification event
type Payload struct {
	Type      string     `json:"type"`
	OccurAt   int64      `json:"occur_at"`
	Operator  string     `json:"operator"`
	EventData *EventData `json:"event_data,omitempty"`
}

// EventData of notification event payload
type EventData struct {
	Resources   []*Resource       `json:"resources,omitempty"`
	Repository  *Repository       `json:"repository,omitempty"`
	Replication *Replication      `json:"replication,omitempty"`
	Custom      map[string]string `json:"custom_attributes,omitempty"`
}

// Resource describe infos of resource triggered notification
type Resource struct {
	Digest       string                 `json:"digest,omitempty"`
	Tag          string                 `json:"tag"`
	ResourceURL  string                 `json:"resource_url,omitempty"`
	ScanOverview map[string]interface{} `json:"scan_overview,omitempty"`
}

// Repository info of notification event
type Repository struct {
	DateCreated  int64  `json:"date_created,omitempty"`
	Name         string `json:"name"`
	Namespace    string `json:"namespace"`
	RepoFullName string `json:"repo_full_name"`
	RepoType     string `json:"repo_type"`
}

// Replication describes replication infos
type Replication struct {
	HarborHostname     string               `json:"harbor_hostname,omitempty"`
	JobStatus          string               `json:"job_status,omitempty"`
	Description        string               `json:"description,omitempty"`
	ArtifactType       string               `json:"artifact_type,omitempty"`
	AuthenticationType string               `json:"authentication_type,omitempty"`
	OverrideMode       bool                 `json:"override_mode,omitempty"`
	TriggerType        string               `json:"trigger_type,omitempty"`
	PolicyCreator      string               `json:"policy_creator,omitempty"`
	ExecutionTimestamp int64                `json:"execution_timestamp,omitempty"`
	SrcResource        *ReplicationResource `json:"src_resource,omitempty"`
	DestResource       *ReplicationResource `json:"dest_resource,omitempty"`
	SuccessfulArtifact []*ArtifactInfo      `json:"successful_artifact,omitempty"`
	FailedArtifact     []*ArtifactInfo      `json:"failed_artifact,omitempty"`
}

// ArtifactInfo describe info of artifact replicated
type ArtifactInfo struct {
	Type       string `json:"type"`
	Status     string `json:"status"`
	NameAndTag string `json:"name_tag"`
	FailReason string `json:"fail_reason,omitempty"`
}

// ReplicationResource describes replication resource info
type ReplicationResource struct {
	RegistryName string `json:"registry_name,omitempty"`
	RegistryType string `json:"registry_type"`
	Endpoint     string `json:"endpoint"`
	Provider     string `json:"provider,omitempty"`
	Namespace    string `json:"namespace,omitempty"`
}
