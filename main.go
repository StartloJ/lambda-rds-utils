package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"reflect"
	"strings"
	"time"
)

type BridgeEvent struct {
	Version    string    `json:"version"`
	ID         string    `json:"id"`
	DetailType string    `json:"detail-type"`
	Source     string    `json:"source"`
	Account    string    `json:"account"`
	Time       time.Time `json:"time"`
	Region     string    `json:"region"`
	Resources  []string  `json:"resources"`
	Detail     struct {
		EventCategories  []string  `json:"EventCategories"`
		SourceType       string    `json:"SourceType"`
		SourceArn        string    `json:"SourceArn"`
		Date             time.Time `json:"Date"`
		SourceIdentifier string    `json:"SourceIdentifier"`
		Message          string    `json:"Message"`
		EventID          string    `json:"EventID"`
	} `json:"detail"`
}

func EnvSetter() {
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("opt")
	viper.AutomaticEnv()
}

// This func input session of RDS service. Then, it get all list of database
// snapshot that already in your resources in List of Snapshot name.
func ListAllSnapshot(svc *rds.RDS) ([]string, error) {

	db := viper.GetString("DB_NAME")

	out, err := svc.DescribeDBSnapshots(&rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: &db,
	})
	if err != nil {
		panic(err)
	}

	var snapShotName []string

	for _, b := range out.DBSnapshots {
		logrus.Infof("We have %s", ConvertToString(*b.DBSnapshotArn))
		snapShotName = append(snapShotName, ConvertToString(*b.DBSnapshotArn))
	}

	return snapShotName, nil
}

// Request to AWS RDS API for create event copy a specific snapshot
// into across region and encrypt it by KMS key multi-region
func CopySnapshotToTarget(svc *rds.RDS, snap string) (string, error) {
	targetSnapArn := strings.Split(snap, ":")
	targetSnapName := targetSnapArn[len(targetSnapArn)-1]

	// Copy the snapshot to the target region
	copySnapshotInput := &rds.CopyDBSnapshotInput{
		OptionGroupName:            aws.String(viper.GetString("option_group_name")),
		KmsKeyId:                   aws.String(viper.GetString("kms_key_id")),
		CopyTags:                   aws.Bool(true),
		SourceDBSnapshotIdentifier: aws.String(snap),
		TargetDBSnapshotIdentifier: aws.String(targetSnapName),
		SourceRegion:               aws.String(viper.GetString("src_region")),
	}

	_, err := svc.CopyDBSnapshot(copySnapshotInput)
	if err != nil {
		logrus.Errorf("Copy request %s is failed", snap)
		return "", err
	}

	logrus.Infof("Copy %s is created", snap)
	return fmt.Sprintf("Copy %s is created", snap), nil
}

func ConvertToString(value interface{}) string {
	if value == nil {
		return ""
	}

	// Check if the value is a pointer and is nil
	if reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil() {
		return ""
	}

	return fmt.Sprintf("%v", value)
}

// this func use to remove duplicate name source if found in target
// you shouldn't switch position input to this func.
// `t` is mean a target that use double check to `s`
// it will remove a value in `s` if found it in `t`
func GetUniqueSnapShots(t, s []string) ([]string, error) {
	//Create a map to keep track a unique strings
	uniqueMap := make(map[string]bool)

	//Iterate over `s` and add each string to map
	for _, str := range s {
		uniqueMap[str] = true
	}

	//Iterate over `t` and remove any string that are already in uniqueMap
	for _, str2 := range t {
		delete(uniqueMap, str2)
	}

	//Convert the unique string from Map to slice string
	result := make([]string, 0, len(uniqueMap))
	for str := range uniqueMap {
		result = append(result, str)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("not any Snapshot unique between source and target region")
	}

	return result, nil
}

func HandlerEvents(event BridgeEvent) error {
	logrus.Infof("we working on event: %s", ConvertToString(event.Detail.EventID))

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(viper.GetString("target_region")),
		},
	}))

	// This condition will processing if the event is DB snapshot created.
	if event.Detail.EventID == "RDS-EVENT-0042" {
		targetSvc := rds.New(sess)
		rep_dbSnapshots, err := ListAllSnapshot(targetSvc)
		if err != nil {
			return fmt.Errorf("we can't get any Snapshot name")
		}
		logrus.Info("gathering DB snapshot from target region is completed")

		SrcSvc := rds.New(sess, aws.NewConfig().WithRegion(viper.GetString("src_region")))
		src_dbSnapShots, err := ListAllSnapshot(SrcSvc)
		if err != nil {
			logrus.Error(err.Error())
		}
		logrus.Info("gathering DB snapshot from source region is completed")

		dbSnapShots2Copy, err := GetUniqueSnapShots(rep_dbSnapshots, src_dbSnapShots)
		if err != nil {
			logrus.Warnf("it doesn't any task copy snapshot to %s", viper.GetString("target_region"))
			return nil
		}
		logrus.Debugf("now, trying to copy %s", strings.Join(dbSnapShots2Copy, ","))

		for s := range dbSnapShots2Copy {
			logrus.Infof("trying to copy DBSnapshot to %s...", viper.GetString("target_region"))
			_, err := CopySnapshotToTarget(targetSvc, dbSnapShots2Copy[s])
			if err != nil {
				logrus.Error(err.Error())
			}
		}
		logrus.Info("all DB snapshot is copying...")

		return nil
	}
	return fmt.Errorf("something else")
}

func Init() {
	logrus.SetReportCaller(viper.GetBool("debug"))
	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetOutput(os.Stdout)
}

func main() {
	EnvSetter()
	Init()
	// Make the hangler available for remote procedure call by Lambda
	logrus.Info("we starting handle lambda...")
	lambda.Start(HandlerEvents)
}

// {
// 	"version": "0",
// 	"id": "844e2571-85d4-695f-b930-0153b71dcb42",
// 	"detail-type": "RDS DB Snapshot Event",
// 	"source": "aws.rds",
// 	"account": "123456789012",
// 	"time": "2018-10-06T12:26:13Z",
// 	"region": "us-east-1",
// 	"resources": ["arn:aws:rds:us-east-1:123456789012:snapshot:rds:snapshot-replica-2018-10-06-12-24"],
// 	"detail": {
// 	  "EventCategories": ["creation"],
// 	  "SourceType": "SNAPSHOT",
// 	  "SourceArn": "arn:aws:rds:us-east-1:123456789012:snapshot:rds:snapshot-replica-2018-10-06-12-24",
// 	  "Date": "2018-10-06T12:26:13.882Z",
// 	  "SourceIdentifier": "rds:snapshot-replica-2018-10-06-12-24",
// 	  "Message": "Automated snapshot created",
// 	  "EventID": "RDS-EVENT-0091"
// 	}
// }
