const { createElement: e, memo, useEffect, useState } = React;

import graphqlQuery from '../graphqlQuery.js';
import ImagePanel, { imageTags } from '../components/ImagePanel.js';

const imagesQuery = `
query($allowTags: [String], $omitTags: [String], $onlyTags: [String]) {
  Images(limit: 12, allowTags: $allowTags, omitTags: $omitTags, onlyTags: $onlyTags) {
    id
    url
    tags
  }
}
`;

export default function Random() {
  const [images, setImages] = useState([]);
  const [loading, setLoading] = useState(false);
  const [firstLoad, setFirstLoad] = useState(true);
  const [variables, setVariables] = useState({});

  const makeOnChangeTagSelector = (variable) => (
    (e) => {
      if (e.target.value === '') {
        setVariables({
          ...variables,
          [variable]: null
        })
      } else {
        setVariables({
          ...variables,
          [variable]: Array.from(e.target.selectedOptions, option => option.value)
        })
      }
    }
  )

  const TagSelector = ({ label, variable }) => (
    e('div', {className: 'd-flex col-2 justify-content-center align-items-center'}, [
      e('label', {key: 'label', className: 'col-form-label flex-grow-1 text-end me-2'}, label),
      e('select', {
        key: 'select',
        multiple: true,
        className: 'form-select w-50',
        onChange: makeOnChangeTagSelector(variable),
        value: variables[variable] || [],
      }, [
        e('option', {key: '__default', value: ''}),
        ...imageTags.map((tag) => e('option', {key: tag, value: tag}, tag))
      ]),
    ])
  )

  const tagSelector = (label, variable) => e(TagSelector, {key: label, label, variable})

  function replaceImage(idx, image) {
    setImages([...images.slice(0, idx), image, ...images.slice(idx+1)])
  }

  function loadImages(reset) {
    setFirstLoad(false);
    setLoading(true);
    graphqlQuery({query: imagesQuery, variables})
      .then(res => res.json())
      .then((result) => {
        if ('errors' in result) {
          console.log(result);
        } else {
          if (reset) {
            setImages(result.data.Images);
          } else {
            setImages([...images, ...result.data.Images]);
          }
        }
        setLoading(false);
      }, (error) => {
        console.log(error);
      });
  }

  function clearAndLoadImages() {
    setImages([]);
    loadImages(true);
  }

  let loadingZone = null;
  if (loading) {
    loadingZone = e('div', {
      key: 'loading',
      className: 'alert alert-primary'
    }, 'Loading...');
  } else if (!firstLoad) {
    loadingZone = e('button', {
      key: 'loading',
      type: 'button',
      className: 'btn btn-primary',
      onClick: () => loadImages()
    }, 'Load More');
  }

  return e('div', {className: 'container-fluid'}, [
    e('div', {key: 'options', className: 'row'}, [
      tagSelector('Allow Tags', 'allowTags'),
      tagSelector('Omit Tags', 'omitTags'),
      tagSelector('Only Tags', 'onlyTags'),
      e('div', {
        key: 'clearimages',
        className: 'col-1 d-flex justify-content-center align-items-center'
      },
        e('button', {
          type: 'button',
          className: 'btn btn-info',
          onClick: () => clearAndLoadImages()
        }, firstLoad ? 'Load Images' : 'Clear And Load New Images'),
      ),
    ]),
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
