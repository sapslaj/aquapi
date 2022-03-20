import os
import praw
import requests
import boto3
import hashlib
import mimetypes

from helpers import *

def generate_tags(post):
    tags = set()
    if post.over_18:
        tags.add("nsfw")
        tags.add("ecchi")
    if post.link_flair_text == "Meme":
        tags.add("meme")
    return tags

if __name__ == "__main__":
    s3 = boto3.client("s3")
    reddit = praw.Reddit(
        client_id=CLIENT_ID, client_secret=CLIENT_SECRET, user_agent=USER_AGENT
    )
    aqua_subreddit = reddit.subreddit("AquaSama")
    for post in aqua_subreddit.top("all", limit=1000):
        try:
            download = requests.get(post.url)
            content_type = download.headers["Content-Type"]
            if content_type not in ("image/jpeg", "image/png"):
                continue
            md5 = hashlib.md5(download.content).hexdigest()
            extension = mimetypes.guess_extension(content_type)
            key = f"{md5}{extension}"
            print(key, end="\t")
            tags = generate_tags(post)
            try:
                object_tagging = s3.get_object_tagging(
                    Bucket=IMAGES_BUCKET,
                    Key=key,
                )
                object_tags = from_tag_set(object_tagging["TagSet"])
                print(f"exists with tags {object_tags}", end="\t")
                tags.update(object_tags)
            except s3.exceptions.NoSuchKey:
                print("does not exist, creating", end="\t")
                s3.put_object(
                    Bucket=IMAGES_BUCKET,
                    ACL="public-read",
                    Body=download.content,
                    ContentType=content_type,
                    Key=key,
                )
            if tags:
                print(f"new tags: {list(tags)}", end="")
                s3.put_object_tagging(
                    Bucket=IMAGES_BUCKET,
                    Key=key,
                    Tagging={"TagSet": to_tag_set(tags)}
                )
            print()
        except Exception as e:
            print(e)
            continue
