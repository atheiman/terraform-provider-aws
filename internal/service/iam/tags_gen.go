// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package iam

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
)

// listTags_Func is the type of the listTags_ function.
type listTags_Func func(context.Context, any, string) error

// updateTags_Func is the type of the updateTags_ function.
type updateTags_Func func(context.Context, any, string, any, any) error

var listTags_ listTags_Func

// []*SERVICE.Tag handling

// Tags returns iam service tags.
func Tags(tags tftags.KeyValueTags) []*iam.Tag {
	result := make([]*iam.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &iam.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from iam service tags.
func KeyValueTags(ctx context.Context, tags []*iam.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsIn returns iam service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) []*iam.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets iam service tags in Context.
func setTagsOut(ctx context.Context, tags []*iam.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = types.Some(KeyValueTags(ctx, tags))
	}
}

var updateTags_ updateTags_Func
