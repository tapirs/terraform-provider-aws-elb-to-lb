// +build ignore

package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

const filename = `update_tags_gen.go`

var serviceNames = []string{
	"accessanalyzer",
	"acm",
	"acmpca",
	"amplify",
	"apigateway",
	"apigatewayv2",
	"appmesh",
	"appstream",
	"appsync",
	"athena",
	"backup",
	"cloud9",
	"cloudfront",
	"cloudhsmv2",
	"cloudtrail",
	"cloudwatch",
	"cloudwatchevents",
	"cloudwatchlogs",
	"codecommit",
	"codedeploy",
	"codepipeline",
	"codestarnotifications",
	"cognitoidentity",
	"cognitoidentityprovider",
	"configservice",
	"databasemigrationservice",
	"dataexchange",
	"datapipeline",
	"datasync",
	"dax",
	"devicefarm",
	"directconnect",
	"directoryservice",
	"dlm",
	"docdb",
	"dynamodb",
	"ec2",
	"ecr",
	"ecs",
	"efs",
	"eks",
	"elasticache",
	"elasticbeanstalk",
	"elasticsearchservice",
	"elb",
	"elbv2",
	"emr",
	"firehose",
	"fsx",
	"gamelift",
	"glacier",
	"globalaccelerator",
	"glue",
	"guardduty",
	"greengrass",
	"imagebuilder",
	"iot",
	"iotanalytics",
	"iotevents",
	"kafka",
	"kinesis",
	"kinesisanalytics",
	"kinesisanalyticsv2",
	"kinesisvideo",
	"kms",
	"lambda",
	"licensemanager",
	"lightsail",
	"mediaconnect",
	"mediaconvert",
	"medialive",
	"mediapackage",
	"mediastore",
	"mq",
	"neptune",
	"opsworks",
	"organizations",
	"pinpoint",
	"qldb",
	"quicksight",
	"ram",
	"rds",
	"redshift",
	"resourcegroups",
	"route53",
	"route53resolver",
	"sagemaker",
	"secretsmanager",
	"securityhub",
	"sfn",
	"sns",
	"sqs",
	"ssm",
	"storagegateway",
	"swf",
	"transfer",
	"waf",
	"wafregional",
	"wafv2",
}

type TemplateData struct {
	ServiceNames []string
}

func main() {
	// Always sort to reduce any potential generation churn
	sort.Strings(serviceNames)

	templateData := TemplateData{
		ServiceNames: serviceNames,
	}
	templateFuncMap := template.FuncMap{
		"ClientType":                      keyvaluetags.ServiceClientType,
		"TagFunction":                     ServiceTagFunction,
		"TagFunctionBatchSize":            ServiceTagFunctionBatchSize,
		"TagInputCustomValue":             ServiceTagInputCustomValue,
		"TagInputIdentifierField":         ServiceTagInputIdentifierField,
		"TagInputIdentifierRequiresSlice": ServiceTagInputIdentifierRequiresSlice,
		"TagInputResourceTypeField":       ServiceTagInputResourceTypeField,
		"TagInputTagsField":               ServiceTagInputTagsField,
		"TagPackage":                      keyvaluetags.ServiceTagPackage,
		"Title":                           strings.Title,
		"UntagFunction":                   ServiceUntagFunction,
		"UntagInputCustomValue":           ServiceUntagInputCustomValue,
		"UntagInputRequiresTagKeyType":    ServiceUntagInputRequiresTagKeyType,
		"UntagInputRequiresTagType":       ServiceUntagInputRequiresTagType,
		"UntagInputTagsField":             ServiceUntagInputTagsField,
	}

	tmpl, err := template.New("updatetags").Funcs(templateFuncMap).Parse(templateBody)

	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)

	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	generatedFileContents, err := format.Source(buffer.Bytes())

	if err != nil {
		log.Fatalf("error formatting generated file: %s", err)
	}

	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("error creating file (%s): %s", filename, err)
	}

	defer f.Close()

	_, err = f.Write(generatedFileContents)

	if err != nil {
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}
}

var templateBody = `
// Code generated by generators/updatetags/main.go; DO NOT EDIT.

package keyvaluetags

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
{{- range .ServiceNames }}
	"github.com/aws/aws-sdk-go/service/{{ . }}"
{{- end }}
)
{{ range .ServiceNames }}

// {{ . | Title }}UpdateTags updates {{ . }} service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func {{ . | Title }}UpdateTags(conn {{ . | ClientType }}, identifier string{{ if . | TagInputResourceTypeField }}, resourceType string{{ end }}, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := New(oldTagsMap)
	newTags := New(newTagsMap)
	{{- if eq (. | TagFunction) (. | UntagFunction) }}
	removedTags := oldTags.Removed(newTags)
	updatedTags := oldTags.Updated(newTags)

	// Ensure we do not send empty requests
	if len(removedTags) == 0 && len(updatedTags) == 0 {
		return nil
	}

	input := &{{ . | TagPackage }}.{{ . | TagFunction }}Input{
		{{- if . | TagInputIdentifierRequiresSlice }}
		{{ . | TagInputIdentifierField }}:   aws.StringSlice([]string{identifier}),
		{{- else }}
		{{ . | TagInputIdentifierField }}:   aws.String(identifier),
		{{- end }}
		{{- if . | TagInputResourceTypeField }}
		{{ . | TagInputResourceTypeField }}: aws.String(resourceType),
		{{- end }}
	}

	if len(updatedTags) > 0 {
		input.{{ . | TagInputTagsField }} = updatedTags.IgnoreAws().{{ . | Title }}Tags()
	}

	if len(removedTags) > 0 {
		{{- if . | UntagInputRequiresTagType }}
		input.{{ . | UntagInputTagsField }} = removedTags.IgnoreAws().{{ . | Title }}Tags()
		{{- else if . | UntagInputRequiresTagKeyType }}
		input.{{ . | UntagInputTagsField }} = removedTags.IgnoreAws().{{ . | Title }}TagKeys()
		{{- else if . | UntagInputCustomValue }}
		input.{{ . | UntagInputTagsField }} = {{ . | UntagInputCustomValue }}
		{{- else }}
		input.{{ . | UntagInputTagsField }} = aws.StringSlice(removedTags.Keys())
		{{- end }}
	}

	_, err := conn.{{ . | TagFunction }}(input)

	if err != nil {
		return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
	}

	{{- else }}

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		{{- if . | TagFunctionBatchSize }}
		for _, removedTags := range removedTags.Chunks({{ . | TagFunctionBatchSize }}) {
		{{- end }}
		input := &{{ . | TagPackage }}.{{ . | UntagFunction }}Input{
			{{- if . | TagInputIdentifierRequiresSlice }}
			{{ . | TagInputIdentifierField }}:   aws.StringSlice([]string{identifier}),
			{{- else }}
			{{ . | TagInputIdentifierField }}:   aws.String(identifier),
			{{- end }}
			{{- if . | TagInputResourceTypeField }}
			{{ . | TagInputResourceTypeField }}: aws.String(resourceType),
			{{- end }}
			{{- if . | UntagInputRequiresTagType }}
			{{ . | UntagInputTagsField }}:       removedTags.IgnoreAws().{{ . | Title }}Tags(),
			{{- else if . | UntagInputRequiresTagKeyType }}
			{{ . | UntagInputTagsField }}:       removedTags.IgnoreAws().{{ . | Title }}TagKeys(),
			{{- else if . | UntagInputCustomValue }}
			{{ . | UntagInputTagsField }}:       {{ . | UntagInputCustomValue }},
			{{- else }}
			{{ . | UntagInputTagsField }}:       aws.StringSlice(removedTags.IgnoreAws().Keys()),
			{{- end }}
		}

		_, err := conn.{{ . | UntagFunction }}(input)

		if err != nil {
			return fmt.Errorf("error untagging resource (%s): %w", identifier, err)
		}
		{{- if . | TagFunctionBatchSize }}
		}
		{{- end }}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		{{- if . | TagFunctionBatchSize }}
		for _, updatedTags := range updatedTags.Chunks({{ . | TagFunctionBatchSize }}) {
		{{- end }}
		input := &{{ . | TagPackage }}.{{ . | TagFunction }}Input{
			{{- if . | TagInputIdentifierRequiresSlice }}
			{{ . | TagInputIdentifierField }}:   aws.StringSlice([]string{identifier}),
			{{- else }}
			{{ . | TagInputIdentifierField }}:   aws.String(identifier),
			{{- end }}
			{{- if . | TagInputResourceTypeField }}
			{{ . | TagInputResourceTypeField }}: aws.String(resourceType),
			{{- end }}
			{{- if . | TagInputCustomValue }}
			{{ . | TagInputTagsField }}:         {{ . | TagInputCustomValue }},
			{{- else }}
			{{ . | TagInputTagsField }}:         updatedTags.IgnoreAws().{{ . | Title }}Tags(),
			{{- end }}
		}

		_, err := conn.{{ . | TagFunction }}(input)

		if err != nil {
			return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
		}
		{{- if . | TagFunctionBatchSize }}
		}
		{{- end }}
	}

	{{- end }}

	return nil
}
{{- end }}
`

// ServiceTagFunction determines the service tagging function.
func ServiceTagFunction(serviceName string) string {
	switch serviceName {
	case "acm":
		return "AddTagsToCertificate"
	case "acmpca":
		return "TagCertificateAuthority"
	case "cloudtrail":
		return "AddTags"
	case "cloudwatchlogs":
		return "TagLogGroup"
	case "databasemigrationservice":
		return "AddTagsToResource"
	case "datapipeline":
		return "AddTags"
	case "directoryservice":
		return "AddTagsToResource"
	case "docdb":
		return "AddTagsToResource"
	case "ec2":
		return "CreateTags"
	case "elasticache":
		return "AddTagsToResource"
	case "elasticbeanstalk":
		return "UpdateTagsForResource"
	case "elasticsearchservice":
		return "AddTags"
	case "elb":
		return "AddTags"
	case "elbv2":
		return "AddTags"
	case "emr":
		return "AddTags"
	case "firehose":
		return "TagDeliveryStream"
	case "glacier":
		return "AddTagsToVault"
	case "kinesis":
		return "AddTagsToStream"
	case "kinesisvideo":
		return "TagStream"
	case "medialive":
		return "CreateTags"
	case "mq":
		return "CreateTags"
	case "neptune":
		return "AddTagsToResource"
	case "rds":
		return "AddTagsToResource"
	case "redshift":
		return "CreateTags"
	case "resourcegroups":
		return "Tag"
	case "route53":
		return "ChangeTagsForResource"
	case "sagemaker":
		return "AddTags"
	case "sqs":
		return "TagQueue"
	case "ssm":
		return "AddTagsToResource"
	case "storagegateway":
		return "AddTagsToResource"
	default:
		return "TagResource"
	}
}

// ServiceTagFunctionBatchSize determines the batch size (if any) for tagging and untagging.
func ServiceTagFunctionBatchSize(serviceName string) string {
	switch serviceName {
	case "kinesis":
		return "10"
	default:
		return ""
	}
}

// ServiceTagInputIdentifierField determines the service tag identifier field.
func ServiceTagInputIdentifierField(serviceName string) string {
	switch serviceName {
	case "acm":
		return "CertificateArn"
	case "acmpca":
		return "CertificateAuthorityArn"
	case "athena":
		return "ResourceARN"
	case "cloud9":
		return "ResourceARN"
	case "cloudfront":
		return "Resource"
	case "cloudhsmv2":
		return "ResourceId"
	case "cloudtrail":
		return "ResourceId"
	case "cloudwatch":
		return "ResourceARN"
	case "cloudwatchevents":
		return "ResourceARN"
	case "cloudwatchlogs":
		return "LogGroupName"
	case "codestarnotifications":
		return "Arn"
	case "datapipeline":
		return "PipelineId"
	case "dax":
		return "ResourceName"
	case "devicefarm":
		return "ResourceARN"
	case "directoryservice":
		return "ResourceId"
	case "docdb":
		return "ResourceName"
	case "ec2":
		return "Resources"
	case "efs":
		return "ResourceId"
	case "elasticache":
		return "ResourceName"
	case "elasticsearchservice":
		return "ARN"
	case "elb":
		return "LoadBalancerNames"
	case "elbv2":
		return "ResourceArns"
	case "emr":
		return "ResourceId"
	case "firehose":
		return "DeliveryStreamName"
	case "fsx":
		return "ResourceARN"
	case "gamelift":
		return "ResourceARN"
	case "glacier":
		return "VaultName"
	case "kinesis":
		return "StreamName"
	case "kinesisanalytics":
		return "ResourceARN"
	case "kinesisanalyticsv2":
		return "ResourceARN"
	case "kinesisvideo":
		return "StreamARN"
	case "kms":
		return "KeyId"
	case "lambda":
		return "Resource"
	case "lightsail":
		return "ResourceName"
	case "mediaconvert":
		return "Arn"
	case "mediastore":
		return "Resource"
	case "neptune":
		return "ResourceName"
	case "organizations":
		return "ResourceId"
	case "ram":
		return "ResourceShareArn"
	case "rds":
		return "ResourceName"
	case "redshift":
		return "ResourceName"
	case "resourcegroups":
		return "Arn"
	case "route53":
		return "ResourceId"
	case "secretsmanager":
		return "SecretId"
	case "sqs":
		return "QueueUrl"
	case "ssm":
		return "ResourceId"
	case "storagegateway":
		return "ResourceARN"
	case "transfer":
		return "Arn"
	case "waf":
		return "ResourceARN"
	case "wafregional":
		return "ResourceARN"
	case "wafv2":
		return "ResourceARN"
	default:
		return "ResourceArn"
	}
}

// ServiceTagInputIdentifierRequiresSlice determines if the service tagging resource field requires a slice.
func ServiceTagInputIdentifierRequiresSlice(serviceName string) string {
	switch serviceName {
	case "ec2":
		return "yes"
	case "elb":
		return "yes"
	case "elbv2":
		return "yes"
	default:
		return ""
	}
}

// ServiceTagInputTagsField determines the service tagging tags field.
func ServiceTagInputTagsField(serviceName string) string {
	switch serviceName {
	case "cloudhsmv2":
		return "TagList"
	case "cloudtrail":
		return "TagsList"
	case "elasticbeanstalk":
		return "TagsToAdd"
	case "elasticsearchservice":
		return "TagList"
	case "glue":
		return "TagsToAdd"
	case "pinpoint":
		return "TagsModel"
	case "route53":
		return "AddTags"
	default:
		return "Tags"
	}
}

// ServiceTagInputCustomValue determines any custom value for the service tagging tags field.
func ServiceTagInputCustomValue(serviceName string) string {
	switch serviceName {
	case "cloudfront":
		return "&cloudfront.Tags{Items: updatedTags.IgnoreAws().CloudfrontTags()}"
	case "kinesis":
		return "aws.StringMap(updatedTags.IgnoreAws().Map())"
	case "pinpoint":
		return "&pinpoint.TagsModel{Tags: updatedTags.IgnoreAws().PinpointTags()}"
	default:
		return ""
	}
}

// ServiceTagInputResourceTypeField determines the service tagging resource type field.
func ServiceTagInputResourceTypeField(serviceName string) string {
	switch serviceName {
	case "route53":
		return "ResourceType"
	case "ssm":
		return "ResourceType"
	default:
		return ""
	}
}

// ServiceUntagFunction determines the service untagging function.
func ServiceUntagFunction(serviceName string) string {
	switch serviceName {
	case "acm":
		return "RemoveTagsFromCertificate"
	case "acmpca":
		return "UntagCertificateAuthority"
	case "cloudtrail":
		return "RemoveTags"
	case "cloudwatchlogs":
		return "UntagLogGroup"
	case "databasemigrationservice":
		return "RemoveTagsFromResource"
	case "datapipeline":
		return "RemoveTags"
	case "directoryservice":
		return "RemoveTagsFromResource"
	case "docdb":
		return "RemoveTagsFromResource"
	case "ec2":
		return "DeleteTags"
	case "elasticache":
		return "RemoveTagsFromResource"
	case "elasticbeanstalk":
		return "UpdateTagsForResource"
	case "elasticsearchservice":
		return "RemoveTags"
	case "elb":
		return "RemoveTags"
	case "elbv2":
		return "RemoveTags"
	case "emr":
		return "RemoveTags"
	case "firehose":
		return "UntagDeliveryStream"
	case "glacier":
		return "RemoveTagsFromVault"
	case "kinesis":
		return "RemoveTagsFromStream"
	case "kinesisvideo":
		return "UntagStream"
	case "medialive":
		return "DeleteTags"
	case "mq":
		return "DeleteTags"
	case "neptune":
		return "RemoveTagsFromResource"
	case "rds":
		return "RemoveTagsFromResource"
	case "redshift":
		return "DeleteTags"
	case "resourcegroups":
		return "Untag"
	case "route53":
		return "ChangeTagsForResource"
	case "sagemaker":
		return "DeleteTags"
	case "sqs":
		return "UntagQueue"
	case "ssm":
		return "RemoveTagsFromResource"
	case "storagegateway":
		return "RemoveTagsFromResource"
	default:
		return "UntagResource"
	}
}

// ServiceUntagInputRequiresTagType determines if the service untagging requires full Tag type.
func ServiceUntagInputRequiresTagType(serviceName string) string {
	switch serviceName {
	case "acm":
		return "yes"
	case "acmpca":
		return "yes"
	case "cloudtrail":
		return "yes"
	case "ec2":
		return "yes"
	default:
		return ""
	}
}

// ServiceUntagInputRequiresTagKeyType determines if a special type for the untagging function tag key field is needed.
func ServiceUntagInputRequiresTagKeyType(serviceName string) string {
	switch serviceName {
	case "elb":
		return "yes"
	default:
		return ""
	}
}

// ServiceUntagInputTagsField determines the service untagging tags field.
func ServiceUntagInputTagsField(serviceName string) string {
	switch serviceName {
	case "acm":
		return "Tags"
	case "acmpca":
		return "Tags"
	case "backup":
		return "TagKeyList"
	case "cloudhsmv2":
		return "TagKeyList"
	case "cloudtrail":
		return "TagsList"
	case "cloudwatchlogs":
		return "Tags"
	case "datasync":
		return "Keys"
	case "ec2":
		return "Tags"
	case "elasticbeanstalk":
		return "TagsToRemove"
	case "elb":
		return "Tags"
	case "glue":
		return "TagsToRemove"
	case "kinesisvideo":
		return "TagKeyList"
	case "resourcegroups":
		return "Keys"
	case "route53":
		return "RemoveTagKeys"
	default:
		return "TagKeys"
	}
}

// ServiceUntagInputCustomValue determines any custom value for the service untagging tags field.
func ServiceUntagInputCustomValue(serviceName string) string {
	switch serviceName {
	case "cloudfront":
		return "&cloudfront.TagKeys{Items: aws.StringSlice(removedTags.IgnoreAws().Keys())}"
	default:
		return ""
	}
}