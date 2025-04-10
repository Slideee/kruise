#!/bin/bash -eu
 
set -o nounset
set -o pipefail
set -o errexit
set -x

go clean --modcache
go mod tidy
go mod vendor

rm -r $SRC/kruise/vendor
go get github.com/AdamKorcz/go-118-fuzz-build/testing


#compile_native_go_fuzzer $SRC/kruise/pkg/controller/workloadspread FuzzPatchFavoriteSubsetMetadataToPod fuzz_patch_favorite_subset_metadata_to_pod
#compile_native_go_fuzzer $SRC/kruise/pkg/controller/workloadspread FuzzPodPreferredScore fuzz_pod_preferred_score
#
#compile_native_go_fuzzer $SRC/kruise/pkg/util/workloadspread FuzzInjectWorkloadSpreadIntoPod fuzz_inject_workloadspread_into_pod
#compile_native_go_fuzzer $SRC/kruise/pkg/util/workloadspread FuzzNestedField fuzz_nested_field
#compile_native_go_fuzzer $SRC/kruise/pkg/util/workloadspread FuzzIsPodSelected fuzz_is_pod_selected
#compile_native_go_fuzzer $SRC/kruise/pkg/util/workloadspread FuzzHasPercentSubset fuzz_has_percent_subset
#
#compile_native_go_fuzzer $SRC/kruise/pkg/webhook/workloadspread/validating FuzzValidateWorkloadSpreadSpec fuzz_validate_workloadspread_spec
#compile_native_go_fuzzer $SRC/kruise/pkg/webhook/workloadspread/validating FuzzValidateWorkloadSpreadConflict fuzz_validate_workloadspread_conflict
#compile_native_go_fuzzer $SRC/kruise/pkg/webhook/workloadspread/validating FuzzValidateWorkloadSpreadTargetRefUpdate fuzz_validate_workloadspread_target_ref_update
#
#compile_native_go_fuzzer $SRC/kruise/pkg/controller/uniteddeployment FuzzParseSubsetReplicas fuzz_parse_subset_replicas
#compile_native_go_fuzzer $SRC/kruise/pkg/controller/uniteddeployment FuzzApplySubsetTemplate fuzz_apply_subset_template
#compile_native_go_fuzzer $SRC/kruise/pkg/controller/uniteddeployment FuzzReplicaAllocator fuzz_replica_allocator
#compile_native_go_fuzzer $SRC/kruise/pkg/webhook/uniteddeployment/validating FuzzValidateUnitedDeploymentSpec fuzz_validate_uniteddeployment_spec


compile_native_go_fuzzer $SRC/kruise/pkg/webhook/resourcedistribution/validating FuzzDeserializeResource fuzz_deserialize_resource
compile_native_go_fuzzer $SRC/kruise/pkg/webhook/resourcedistribution/validating FuzzValidateResourceDistributionSpec fuzz_validate_resource_distribution_spec
compile_native_go_fuzzer $SRC/kruise/pkg/webhook/resourcedistribution/validating FuzzValidateResourceDistributionTargets fuzz_validate_resource_distribution_targets
compile_native_go_fuzzer $SRC/kruise/pkg/webhook/resourcedistribution/validating FuzzValidateResourceDistributionResource fuzz_validate_resource_distribution_resource
