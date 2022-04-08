# probably should only run this if you know what you're doing...

import json
import os

from flask import Flask, request

from db import Image, Session

app = Flask(__name__)

def paginate(query, page, page_size):
    return query.limit(page_size).offset((page - 1) * page_size).all()

def serialize_image(image):
    return f"https://aquapic.sapslaj.com/{image.key}"

@app.route("/images")
def images():
    page = int(request.args.get('page', '1'))
    page_size = int(request.args.get('page_size', '50'))
    with Session() as session:
        items = paginate(session.query(Image), page, page_size)
    return {
        "page": page,
        "page_size": page_size,
        "items": list(map(serialize_image, items)),
    }

@app.route("/")
def index():
    return """
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Wreck your browser!</title>
  <style type="text/css">
    .container {
      display: grid;
      grid-template-columns: repeat(8, 1fr);
      grid-template-rows: masonry;
    }
    img {
      width: 100%;
      height: auto;
    }
  </style>
</head>
<body>
  <div class="container"></div>
  <script>
    const container = document.querySelector(".container");
    let page = 0;
    let loading = false;
    let timeout;
    const moarImages = () => {
      if (loading) {
        return
      }
      loading = true;
      page += 1;
      fetch(`/images?page=${page}&page_size=16`)
        .then(response => response.json())
        .then(data => {
          data.items.forEach((image) => {
            const el = document.createElement('img');
            el.setAttribute('src', image)
            container.append(el);
          });
          window.scrollTo(0, document.body.scrollHeight);
          loading = false;
          if (data.items.length > 0) {
            timeout = setTimeout(moarImages, 10000);
          }
        }).catch((reason) => {
          alert(reason);
        });
    }
    moarImages()
  </script>
</body>
</html>
"""

if __name__ == "__main__":
    os.environ["FLASK_ENV"] = "development"
    app.run(debug=True)
