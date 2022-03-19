const { createElement: e, memo, useEffect, useState } = React;

import graphqlQuery from '../graphqlQuery.js';
import ImagePanel from '../components/ImagePanel.js';

const imagesQuery = `
query {
  Images(limit: 12) {
    id
    url
    tags
  }
}
`;

export default function Random() {
  const [images, setImages] = useState([]);
  const [loading, setLoading] = useState(true);

  function replaceImage(idx, image) {
    setImages([...images.slice(0, idx), image, ...images.slice(idx+1)])
  }

  function loadImages() {
    setLoading(true);
    graphqlQuery({query: imagesQuery})
      .then(res => res.json())
      .then((result) => {
        if ('errors' in result) {
          console.log(result);
        } else {
          setImages([...images, ...result.data.Images]);
        }
        setLoading(false);
      }, (error) => {
        console.log(error);
      });
  }

  useEffect(() => {
    if (images.length === 0) {
      loadImages()
    }
  })

  let loadingZone = null;
  if (loading) {
    loadingZone = e('div', {
      key: 'loading',
      className: 'alert alert-primary'
    }, 'Loading...');
  } else {
    loadingZone = e('button', {
      key: 'loading',
      type: 'button',
      className: 'btn btn-primary',
      onClick: () => loadImages()
    }, 'Load More');
  }

  return e('div', {className: 'container-fluid'}, [
    e('div', {
      key: 'row',
      className: 'row row-cols-4'
    }, images.map((i, idx) => e(memo(ImagePanel), {
      ...i,
      idx,
      replaceImage,
      key: idx
    }))),
    e('div', {
      key: 'loading-row',
      className: 'row justify-content-md-center'
    }, e('div', {
      key: 'loading-col',
      className: 'col-2'
    }, loadingZone)),
  ])
}
