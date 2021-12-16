package exporters

import (
	"bufio"
	"context"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

// S3 is a `Storer` interface that puts an ExportResult to the specified S3 bucket. Don't use your main AWS keys for this!! Create read-only keys using IAM
type S3 struct {
	Endpoint     string
	Region       string
	Bucket       string
	AccessKey    string
	ClientSecret string
	UseSSL       bool
}

// Store puts an `ExportResult` struct to an S3 bucket within the specified directory
func (x *S3) Store(result *ExportResult, directory string) *Error {
	if result.Error != nil {
		return result.Error
	}

	file, err := os.Open(result.Path)
	if err != nil {
		return makeErr(err, "")
	}
	defer file.Close()

	buffy := bufio.NewReader(file)
	stat, err := file.Stat()
	if err != nil {
		return makeErr(err, "")
	}

	size := stat.Size()

	ctx := context.Background()
	minioClient, err := minio.New(x.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(x.AccessKey, x.ClientSecret, ""),
		Secure: x.UseSSL,
	})
	if err != nil {
		return makeErr(err, "")
	}
	err = minioClient.MakeBucket(
		ctx,
		x.Bucket,
		minio.MakeBucketOptions{Region: x.Region},
	)

	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, x.Bucket)
		if errBucketExists == nil && exists {
			logrus.Printf("We already own %s\n", x.Bucket)
		} else {
			return makeErr(err, "")
		}
	} else {
		logrus.Infof("Successfully created %s\n", x.Bucket)
	}
	info, err := minioClient.PutObject(
		ctx,
		x.Bucket,
		directory+result.Filename(),
		buffy,
		size,
		minio.PutObjectOptions{ContentType: result.MIME})
	if err != nil {
		return makeErr(err, "")
	}
	logrus.Infof("Successfully uploaded %s of size %d\n", directory+result.Filename(), info.Size)
	return makeErr(err, "")
}
