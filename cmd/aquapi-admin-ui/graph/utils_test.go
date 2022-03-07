package graph

import (
	"fmt"
	"testing"

	"github.com/aws/smithy-go/ptr"
)

type __assertApplicableImage struct {
	tags      []string
	allowTags []*string
	omitTags  []*string
	onlyTags  []*string
}

func (a *__assertApplicableImage) withAllowTags(tags []string) *__assertApplicableImage {
	a.allowTags = ptr.StringSlice(tags)
	return a
}

func (a *__assertApplicableImage) withOmitTags(tags []string) *__assertApplicableImage {
	a.omitTags = ptr.StringSlice(tags)
	return a
}

func (a *__assertApplicableImage) withOnlyTags(tags []string) *__assertApplicableImage {
	a.onlyTags = ptr.StringSlice(tags)
	return a
}

func (a *__assertApplicableImage) strrepr() string {
	allowTags := "nil"
	if a.allowTags != nil {
		allowTags = fmt.Sprintf("%v", ptr.ToStringSlice(a.allowTags))
	}
	omitTags := "nil"
	if a.omitTags != nil {
		omitTags = fmt.Sprintf("%v", ptr.ToStringSlice(a.omitTags))
	}
	onlyTags := "nil"
	if a.onlyTags != nil {
		onlyTags = fmt.Sprintf("%v\n", ptr.ToStringSlice(a.onlyTags))
	}
	return fmt.Sprintf(
		"applicableImage(tags=%s, allowTags=%s, omitTags=%s, onlyTags=%s)",
		fmt.Sprintf("%v", a.tags),
		allowTags,
		omitTags,
		onlyTags,
	)
}

func (a *__assertApplicableImage) isTrue(t *testing.T) {
	if !applicableImage(a.tags, a.allowTags, a.omitTags, a.onlyTags) {
		t.Errorf(
			"%s was false",
			a.strrepr(),
		)
	}
}

func (a *__assertApplicableImage) isFalse(t *testing.T) {
	if applicableImage(a.tags, a.allowTags, a.omitTags, a.onlyTags) {
		t.Errorf(
			"%s was true",
			a.strrepr(),
		)
	}
}

func assertApplicableImage(tags []string) *__assertApplicableImage {
	new := __assertApplicableImage{
		tags: tags,
	}
	return &new
}

func TestApplicableImageBare(t *testing.T) {
	assertApplicableImage([]string{}).isTrue(t)
	assertApplicableImage([]string{"nsfw", "hidden"}).isTrue(t)
}

func TestApplicableImageAllowTags(t *testing.T) {
	assertApplicableImage([]string{}).withAllowTags([]string{}).isFalse(t)
	assertApplicableImage([]string{}).withAllowTags([]string{"nsfw"}).isFalse(t)
	assertApplicableImage([]string{"nsfw"}).withAllowTags([]string{"nsfw"}).isTrue(t)
	assertApplicableImage([]string{"nsfw"}).withAllowTags([]string{"hidden"}).isFalse(t)
	assertApplicableImage([]string{"nsfw", "hidden"}).withAllowTags([]string{"nsfw"}).isTrue(t)
	assertApplicableImage([]string{"nsfw", "hidden"}).withAllowTags([]string{"meme", "nsfw"}).isTrue(t)
	assertApplicableImage([]string{"nsfw"}).withAllowTags([]string{"hidden", "nsfw"}).isTrue(t)
}

func TestApplicableImageOmitTags(t *testing.T) {
	assertApplicableImage([]string{}).withOmitTags([]string{}).isTrue(t)
	assertApplicableImage([]string{}).withOmitTags([]string{"nsfw"}).isTrue(t)
	assertApplicableImage([]string{"nsfw"}).withOmitTags([]string{"nsfw"}).isFalse(t)
	assertApplicableImage([]string{"nsfw"}).withOmitTags([]string{"hidden"}).isTrue(t)
	assertApplicableImage([]string{"nsfw", "hidden"}).withOmitTags([]string{"nsfw"}).isFalse(t)
	assertApplicableImage([]string{"nsfw", "hidden"}).withOmitTags([]string{"meme", "nsfw"}).isFalse(t)
	assertApplicableImage([]string{"nsfw"}).withOmitTags([]string{"hidden", "nsfw"}).isFalse(t)
}

func TestApplicableImageAllowTagsOmitTags(t *testing.T) {
	assertApplicableImage([]string{"nsfw", "hidden"}).withAllowTags([]string{"nsfw"}).withOmitTags([]string{"hidden"}).isFalse(t)
}

func TestApplicableImageOnlyTags(t *testing.T) {
	assertApplicableImage([]string{}).withOnlyTags([]string{}).isTrue(t)
	assertApplicableImage([]string{}).withOnlyTags([]string{"nsfw"}).isFalse(t)
	assertApplicableImage([]string{"nsfw"}).withOnlyTags([]string{"nsfw"}).isTrue(t)
	assertApplicableImage([]string{"nsfw", "hidden"}).withOnlyTags([]string{"hidden", "nsfw"}).isTrue(t)
	assertApplicableImage([]string{"nsfw"}).withOnlyTags([]string{"nsfw", "hidden"}).isFalse(t)
	assertApplicableImage([]string{"nsfw", "hidden"}).withOnlyTags([]string{"nsfw"}).isFalse(t)
}
