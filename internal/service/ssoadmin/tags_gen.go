// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package ssoadmin

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssoadmin"
	"github.com/aws/aws-sdk-go/service/ssoadmin/ssoadminiface"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// listTags_Func is the type of the listTags_ function.
type listTags_Func func(context.Context, any, string, string) error

// updateTags_Func is the type of the updateTags_ function.
type updateTags_Func func(context.Context, any, string, string, any, any) error

// listTags lists ssoadmin service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func listTags(ctx context.Context, conn ssoadminiface.SSOAdminAPI, identifier, resourceType string) (tftags.KeyValueTags, error) {
	input := &ssoadmin.ListTagsForResourceInput{
		ResourceArn: aws.String(identifier),
		InstanceArn: aws.String(resourceType),
	}

	output, err := conn.ListTagsForResourceWithContext(ctx, input)

	if err != nil {
		return tftags.New(ctx, nil), err
	}

	return KeyValueTags(ctx, output.Tags), nil
}

// listTags_ lists ssoadmin service tags and set them in Context.
// It is called from outside this package.
var listTags_ listTags_Func = func(ctx context.Context, meta any, identifier, resourceType string) error {
	tags, err := listTags(ctx, meta.(*conns.AWSClient).SSOAdminConn(ctx), identifier, resourceType)

	if err != nil {
		return err
	}

	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = types.Some(tags)
	}

	return nil
}

// []*SERVICE.Tag handling

// Tags returns ssoadmin service tags.
func Tags(tags tftags.KeyValueTags) []*ssoadmin.Tag {
	result := make([]*ssoadmin.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &ssoadmin.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from ssoadmin service tags.
func KeyValueTags(ctx context.Context, tags []*ssoadmin.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsIn returns ssoadmin service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) []*ssoadmin.Tag {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets ssoadmin service tags in Context.
func setTagsOut(ctx context.Context, tags []*ssoadmin.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = types.Some(KeyValueTags(ctx, tags))
	}
}

// updateTags updates ssoadmin service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTags(ctx context.Context, conn ssoadminiface.SSOAdminAPI, identifier, resourceType string, oldTagsMap, newTagsMap any) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.SSOAdmin)
	if len(removedTags) > 0 {
		input := &ssoadmin.UntagResourceInput{
			ResourceArn: aws.String(identifier),
			InstanceArn: aws.String(resourceType),
			TagKeys:     aws.StringSlice(removedTags.Keys()),
		}

		_, err := conn.UntagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.SSOAdmin)
	if len(updatedTags) > 0 {
		input := &ssoadmin.TagResourceInput{
			ResourceArn: aws.String(identifier),
			InstanceArn: aws.String(resourceType),
			Tags:        Tags(updatedTags),
		}

		_, err := conn.TagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

// updateTags_ updates ssoadmin service tags.
// It is called from outside this package.
var updateTags_ updateTags_Func = func(ctx context.Context, meta any, identifier, resourceType string, oldTags, newTags any) error {
	return updateTags(ctx, meta.(*conns.AWSClient).SSOAdminConn(ctx), identifier, resourceType, oldTags, newTags)
}
