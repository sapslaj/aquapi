import os
import praw
import requests
import boto3
import hashlib
import mimetypes

USER_AGENT = "AquAPI Reddit Scaper (by /u/sapslaj)"
CLIENT_ID = os.environ["REDDIT_CLIENT_ID"]
CLIENT_SECRET = os.environ["REDDIT_CLIENT_SECRET"]
IMAGES_BUCKET = os.environ["AQUAPI_IMAGES_BUCKET"]

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
            print(key)
            s3.put_object(
                Bucket=IMAGES_BUCKET,
                ACL="public-read",
                Body=download.content,
                ContentType=content_type,
                Key=key,
            )
        except Exception as e:
            print(e)
            continue
