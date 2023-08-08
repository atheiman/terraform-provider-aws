// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package macie2

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/types"
)

// listTags_Func is the type of the listTags_ function.
type listTags_Func func(context.Context, any, string) error

// updateTags_Func is the type of the updateTags_ function.
type updateTags_Func func(context.Context, any, string, any, any) error

var listTags_ listTags_Func

// map[string]*string handling

// Tags returns macie2 service tags.
func Tags(tags tftags.KeyValueTags) map[string]*string {
	return aws.StringMap(tags.Map())
}

// KeyValueTags creates tftags.KeyValueTags from macie2 service tags.
func KeyValueTags(ctx context.Context, tags map[string]*string) tftags.KeyValueTags {
	return tftags.New(ctx, tags)
}

// getTagsIn returns macie2 service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) map[string]*string {
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags) > 0 {
			return tags
		}
	}

	return nil
}

// setTagsOut sets macie2 service tags in Context.
func setTagsOut(ctx context.Context, tags map[string]*string) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = types.Some(KeyValueTags(ctx, tags))
	}
}

var updateTags_ updateTags_Func
