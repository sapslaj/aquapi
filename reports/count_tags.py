from functools import reduce
import operator

from sqlalchemy import func, and_, or_

from db import Image, Session


if __name__ == "__main__":
    with Session() as session:

        def image_tags_selector(tags, has_tag=True):
            return [getattr(Image, tag) == has_tag for tag in tags]

        def image_filter_count(f):
            return session.query(Image).filter(f).count()

        def image_tags_count(tags, has_tag=True, op=None):
            if not op:
                op = or_ if has_tag else and_
            return image_filter_count(op(*image_tags_selector(tags, has_tag=has_tag)))

        def total_percent(count):
            return f"{(count / total)*100}%"

        def print_count(name, value):
            print(name.rjust(15), value, total_percent(value), sep="\t")

        total = session.query(Image).count()
        print("Total".rjust(15), total, sep="\t")
        print()

        tags = [
            "hidden",
            "nsfw",
            "ecchi",
            "hentai",
            "meme",
            "collage",
        ]
        for tag in tags:
            count = image_tags_count(tags=[tag])
            print_count(tag, count)
        print()

        print_count("has tag", image_tags_count(tags=tags))
        print_count("no tag", image_tags_count(tags=tags, has_tag=False))
        print()

        print_count("is nsfw", image_tags_count(tags=["nsfw", "ecchi", "hentai"]))
        print_count("not nsfw", image_tags_count(tags=["nsfw", "ecchi", "hentai"], has_tag=False))
        print_count(
            "no nsfw",
            image_filter_count(
                and_(
                    *image_tags_selector(tags=["nsfw"], has_tag=False),
                    or_(*image_tags_selector(tags=["ecchi", "hentai"])),
                )
            ),
        )
        print_count(
            "just nsfw",
            image_filter_count(
                and_(
                    *image_tags_selector(tags=["nsfw"], has_tag=True),
                    *image_tags_selector(tags=["ecchi", "hentai"], has_tag=False),
                )
            ),
        )
