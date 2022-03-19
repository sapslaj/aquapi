const { createElement: e } = React;

import ImageTagToggle from './ImageTagToggle.js';

const imageTags = ['hidden', 'nsfw', 'ecchi', 'hentai', 'meme', 'collage'];

export default function ImagePanel({idx, replaceImage, id, url, tags}) {
  const thisImageTags = imageTags.reduce((o, key) => ({ ...o, [key]: tags.includes(key)}), {})

  function updateImage() {
    if (replaceImage) {
      replaceImage(idx, {id, url, tags: imageTags.filter((tag) => thisImageTags[tag])})
    }
  }

  function updateImageTag(tag, value) {
    thisImageTags[tag] = value
    updateImage()
  }

  return e('div', {key: `div-col-${id}`, className: 'col'}, [
    e('div', {
      key: `div-card-${id}`,
      className: 'card',
      style: {
        // hack
        marginBottom: 'var(--bs-gutter-x)'
      }
    }, [
      e('img', {key: `div-card-img-${id}`, className: 'card-img-top', src: url}),
      e('div', {key: `div-card-cody-${id}`, className: 'card-body'}, [
        `id: ${id}`,
        e('br', {key: `div-card-body-br1-${id}`}),
        e('div', {key: `div-card-body-tags-${id}`}, imageTags.map(tag => (
          e(ImageTagToggle, {
            key: `div-card-body-tags-${id}-${tag}`,
            updateImageTag,
            id,
            tag,
            value: thisImageTags[tag]
          })
        )))
      ])
    ])
  ])
}
