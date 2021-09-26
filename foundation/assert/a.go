package a

import (
	"testing"
)

// projectIdFixture = "dev8celbux"
// config           = di.Config{
// 	ProjectID:           "dev8celbux",
// 	CredentialsFilePath: "/service/key.json",
// }

// AssertCalledOnce asserts a given value to 1
func AssertCalledOnce(t *testing.T, got int) {
	t.Helper()
	if got != 1 {
		t.Errorf("Wanted 1 call but got %d calls", got)
	}
}

// AssertEmpty asserts a given string is equal to ""
func AssertEmpty(t *testing.T, val string) {
	t.Helper()
	if val != "" {
		t.Errorf("wanted an empty string but got %v", val)
	}
}

// AssertStrings asserts 2 strings are equal in value
func AssertStrings(t *testing.T, want string, got string) {
	t.Helper()
	if want != got {
		t.Errorf("wanted %v but got %v", want, got)
	}
}

// AssertInts asserts 2 integers are equal in value
func AssertInts(t *testing.T, want int, got int) {
	t.Helper()
	if want != got {
		t.Errorf("wanted %v but got %v", want, got)
	}
}

// AssertInt64s asserts 2 integers are equal in value
func AssertInt64s(t *testing.T, want int64, got int64) {
	t.Helper()
	if want != got {
		t.Errorf("wanted %v but got %v", want, got)
	}
}

// AssertFailure asserts a given error is NOT nil
func AssertFailure(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("Expected Failure")
	}
}

// AssertSuccess asserts a given error is nil
func AssertSuccess(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Expected success but got: %v", err.Error())
	}
}

// AssertOneCall asserts that a function was only called once,
// usually used for stubbing mocks
func AssertOneCall(t *testing.T, got uint) {
	t.Helper()
	if got != 1 {
		t.Errorf("Expected one call, but got called %d times", got)
	}
}

// AssertNoErr asserts a given error is nil
func AssertNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("Expected nil err but got %s", err.Error())
	}
}

// AssertWorkQueue asserts 2 structs are equal in value
// func AssertWorkQueue(t *testing.T, want *b.WorkQueue, got *b.WorkQueue) {
// 	t.Helper()
// 	if !reflect.DeepEqual(want, got) {
// 		t.Errorf("wanted %v but got %v", want, got)
// 	}
// }

// SetupDisburseTest clears BigQuery & Datastore, and reuploads PRODDisbCheck file
// func SetupDisburseTest(ctx context.Context, dependencies di.Dependencies, file []byte, fileName string) error {
// 	// Clear Datastore
// 	kindsToClear := []string{"WorkQueue", "WorkQueue2", "DisbursementSummary"}
// 	err := ClearDatastore(ctx, dependencies, "", kindsToClear...)
// 	if err != nil {
// 		return err
// 	}

// 	// Clear BigQuery
// 	err = ClearBigQuery(ctx, dependencies, "dev8celbux.bulkvoucherwrites.validationEntries")
// 	if err != nil {
// 		return err
// 	}

// 	// Upload file to Cloud Storage
// 	err = UploadCloudStorageFile(ctx, dependencies, file, fileName)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// SetupConfirmDisburseTest prepares dependencies for a disbursement test
// func SetupConfirmDisburseTest(ctx context.Context, dependencies di.Dependencies, namespace string, uidx string, voucherNo int64) error {
// 	f15Key := datastore.IDKey("F1_15", voucherNo, nil)
// 	t15Key := datastore.IDKey("T1_15", voucherNo, nil)
// 	f15Key.Namespace = namespace
// 	t15Key.Namespace = namespace

// 	// Lives in default namespace at all times
// 	summaryKey := datastore.NameKey("DisbursementSummary", uidx, nil)

// 	voucher := b.F1_15{
// 		Amount:      10,
// 		F8:          "ZAR",
// 		Title:       "0648503047",
// 		Audit:       datastore.IDKey("F1_11", 1, nil),
// 		LinkTo:      datastore.IDKey("T1_15", 1, nil),
// 		Owner:       datastore.IDKey("F1_4", 1, nil),
// 		VoucherType: datastore.IDKey("F1_9", 1, nil),
// 		DateCreated: "2020-02-28",
// 	}

// 	link := b.T1_15{
// 		LinkTo:          datastore.IDKey("T1_15", 1, nil),
// 		LastTransaction: time.Now().UnixNano() / 1000 * time.Millisecond.Milliseconds(),
// 	}

// 	summary := b.DisbursementSummary{
// 		Namespace: namespace,
// 		Status:    "PENDING CONFIRMATION",
// 		VoucherNo: "1",
// 	}

// 	_, err := dependencies.Datastore.Put(ctx, f15Key, &voucher)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = dependencies.Datastore.Put(ctx, t15Key, &link)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = dependencies.Datastore.Put(ctx, summaryKey, &summary)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// SetupRollbackTest prepares dependencies for a rollback
// func SetupRollbackTest(ctx context.Context, dependencies di.Dependencies, uidx string) error {
// 	service := disburse.Service{
// 		MerchantAPI: dependencies.MerchantAPI,
// 		Datastore:   dependencies.Datastore,
// 		Log:         dependencies.Log,
// 	}

// 	err := service.ShardTreasury(ctx, 100, "ZHE0VWN/R3QwP2ZyYV1bbHVd", uidx)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// ClearDatastore clears all records in specified Datastore Kinds. Empty namespace selects the default namespace
// func ClearDatastore(ctx context.Context, dependencies di.Dependencies, namespace string, kindsToClear ...string) error {
// 	// Delete all records in specified Kinds
// 	for _, kind := range kindsToClear {
// 		// Get all keys
// 		q := datastore.NewQuery(kind).KeysOnly().Namespace(namespace)
// 		keys, err := dependencies.Datastore.GetAll(ctx, q, nil)
// 		if err != nil {
// 			return err
// 		}

// 		// Delete all records in chunks of 500 or less
// 		for i := 0; i < len(keys); i += 500 {
// 			chunk := min(len(keys)-i, 500)
// 			err = dependencies.Datastore.DeleteMulti(ctx, keys[i:i+chunk])
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }

// ClearBigQuery clears all records in specified BQ tables
// func ClearBigQuery(ctx context.Context, dependencies di.Dependencies, bqTablesToClear ...string) error {
// 	for _, table := range bqTablesToClear {
// 		// Create query
// 		s := fmt.Sprintf("DELETE FROM %v WHERE true;", table)
// 		q := dependencies.BigQuery.Query(s)

// 		// Run query
// 		job, err := q.Run(ctx)
// 		if err != nil {
// 			return err
// 		}

// 		// Wait for query to complete
// 		_, err = job.Wait(ctx)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// min returns the minimum of 2 numbers
// func min(num1 int, num2 int) int {
// 	if num1 > num2 {
// 		return num2
// 	}
// 	return num1
// }

// UploadCloudStorageFile will upload a file to GCP cloud storage
// func UploadCloudStorageFile(ctx context.Context, dependencies di.Dependencies, file []byte, bucketFileName string) error {
// 	bucket := "dev8disburse"
// 	wc := dependencies.CloudStorage.Bucket(bucket).Object(bucketFileName).NewWriter(ctx)
// 	wc.ContentType = "text/plain"

// 	_, err := wc.Write(file)
// 	if err != nil {
// 		return err
// 	}
// 	err = wc.Close()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
