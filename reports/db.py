import os
import asyncio
import concurrent.futures

import boto3
from sqlalchemy import Column, DateTime, Integer, String, Boolean
from sqlalchemy import create_engine
from sqlalchemy.orm import declarative_base, sessionmaker

TAGS_TAG_KEY = "AquaPITags"
TAG_SEPERATOR = ":"

IMAGES_BUCKET = os.environ["AQUAPI_IMAGES_BUCKET"]

engine = create_engine("sqlite:///db.sqlite3")
Session = sessionmaker(bind=engine)

Model = declarative_base()


def from_tag_set(tag_set):
    kvs = [kv for kv in tag_set if kv["Key"] == TAGS_TAG_KEY]
    if not kvs:
        return []
    return kvs[0]["Value"].split(TAG_SEPERATOR)


class Image(Model):
    __tablename__ = "images"

    key = Column(String, primary_key=True)
    hidden = Column(Boolean, default=False)
    nsfw = Column(Boolean, default=False)
    ecchi = Column(Boolean, default=False)
    hentai = Column(Boolean, default=False)
    meme = Column(Boolean, default=False)
    collage = Column(Boolean, default=False)


Model.metadata.create_all(engine)
s3_client = boto3.client("s3")
s3_resource = boto3.resource("s3")
bucket = s3_resource.Bucket(IMAGES_BUCKET)


def image_factory(key):
    def image():
        with Session() as session:
            if key in ("error.html", "index.html", "favicon.ico"):
                return
            object_tagging = s3_client.get_object_tagging(
                Bucket=IMAGES_BUCKET,
                Key=key,
            )
            tags = from_tag_set(object_tagging["TagSet"])
            image = session.query(Image).filter(Image.key == key).first()
            exists = "exists"
            if not image:
                image = Image(key=key)
                exists = "new"
            for tag in tags:
                setattr(image, tag, True)
            session.add(image)
            session.commit()
            print(key, exists, tags, sep="\t")

    return image


async def main(executor):
    event_loop = asyncio.get_event_loop()
    completed, pending = await asyncio.wait(
        [
            event_loop.run_in_executor(executor, image_factory(object_summary.key))
            for object_summary in bucket.objects.all()
        ]
    )
    results = [t.result() for t in completed]


if __name__ == "__main__":
    executor = concurrent.futures.ThreadPoolExecutor(max_workers=64)
    event_loop = asyncio.get_event_loop()
    event_loop.run_until_complete(main(executor))
