package txnlog

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

// CloudWatchClient defines the interface for writing logs to CloudWatch.
// This allows for mocking in tests.
type CloudWatchClient interface {
	WriteLog(entry *LogEntry) error
}

// cloudWatchClient is the production implementation using AWS SDK.
type cloudWatchClient struct {
	client       *cloudwatchlogs.CloudWatchLogs
	logGroup     string
	logStream    string
	sequenceToken *string
}

// NewCloudWatchClient creates a new CloudWatch client.
// If sess is nil, it attempts to create a default session.
func NewCloudWatchClient(sess *session.Session, logGroup, logStream string) (CloudWatchClient, error) {
	if sess == nil {
		var err error
		sess, err = session.NewSession()
		if err != nil {
			return nil, fmt.Errorf("failed to create AWS session: %w", err)
		}
	}

	client := cloudwatchlogs.New(sess)

	cwc := &cloudWatchClient{
		client:    client,
		logGroup:  logGroup,
		logStream: logStream,
	}

	// Ensure log stream exists
	if err := cwc.ensureLogStream(); err != nil {
		return nil, err
	}

	return cwc, nil
}

// ensureLogStream creates the log stream if it doesn't exist.
func (c *cloudWatchClient) ensureLogStream() error {
	// Try to create the log stream; ignore error if it already exists
	_, err := c.client.CreateLogStream(&cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(c.logGroup),
		LogStreamName: aws.String(c.logStream),
	})
	if err != nil {
		// ResourceAlreadyExistsException is fine
		if _, ok := err.(*cloudwatchlogs.ResourceAlreadyExistsException); !ok {
			return fmt.Errorf("failed to create log stream: %w", err)
		}
	}
	return nil
}

// WriteLog writes a log entry to CloudWatch.
func (c *cloudWatchClient) WriteLog(entry *LogEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}

	input := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(c.logGroup),
		LogStreamName: aws.String(c.logStream),
		LogEvents: []*cloudwatchlogs.InputLogEvent{
			{
				Message:   aws.String(string(data)),
				Timestamp: aws.Int64(entry.Timestamp.UnixNano() / int64(time.Millisecond)),
			},
		},
	}

	if c.sequenceToken != nil {
		input.SequenceToken = c.sequenceToken
	}

	output, err := c.client.PutLogEvents(input)
	if err != nil {
		return fmt.Errorf("failed to put log events: %w", err)
	}

	c.sequenceToken = output.NextSequenceToken
	return nil
}

// localFallbackClient logs to stdout when CloudWatch is unavailable.
type localFallbackClient struct{}

// NewLocalFallbackClient creates a client that writes to local log.
func NewLocalFallbackClient() CloudWatchClient {
	return &localFallbackClient{}
}

// WriteLog writes the entry to local log output.
func (c *localFallbackClient) WriteLog(entry *LogEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %w", err)
	}
	log.Printf("[txnlog] %s", string(data))
	return nil
}

// noopClient discards all logs (useful for testing).
type noopClient struct{}

// NewNoopClient creates a client that discards logs.
func NewNoopClient() CloudWatchClient {
	return &noopClient{}
}

// WriteLog does nothing.
func (c *noopClient) WriteLog(entry *LogEntry) error {
	return nil
}
