import os

USER_AGENT = "AquAPI Reddit Scaper (by /u/sapslaj)"
CLIENT_ID = os.environ["REDDIT_CLIENT_ID"]
CLIENT_SECRET = os.environ["REDDIT_CLIENT_SECRET"]
IMAGES_BUCKET = os.environ["AQUAPI_IMAGES_BUCKET"]

TAGS_TAG_KEY = "AquaPITags"
TAG_SEPERATOR = ":"

def from_tag_set(tag_set):
    kvs = [kv for kv in tag_set if kv["Key"] == TAGS_TAG_KEY]
    if not kvs:
        return []
    return kvs[0]["Value"].split(TAG_SEPERATOR)

def to_tag_set(tags):
    return [{"Key": TAGS_TAG_KEY, "Value": TAG_SEPERATOR.join(tags)}]
