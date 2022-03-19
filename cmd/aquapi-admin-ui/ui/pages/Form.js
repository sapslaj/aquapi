import ImagePanel from '../components/ImagePanel.js';
import { debounce } from '../debounce.js'
import graphqlQuery from '../graphqlQuery.js';

const { createElement: e, useState } = React;

const imagesQuery = `
query($key: String!) {
  Image(id: $key) {
    id
    url
    tags
  }
}
`;

const unFuckedKey = (key) => {
  try {
    const url = new URL(key);
    return url.pathname.substring(1);
  } catch (TypeError) {
    return key;
  }
}

function ImagePanelForm({ loaded }) {
  const [image, setImage] = useState(null);

  const updateImage = (key) => {
    graphqlQuery({query: imagesQuery, variables: {key}})
      .then(res => res.json())
      .then((result) => {
        if ('errors' in result) {
          console.log(result);
        } else {
          if (!image) {
            loaded()
          }
          setImage(result.data.Image);
        }
      }, (error) => {
        console.log(error);
      });
  }

  const onChange = (e) => {
    debounce(() => {
      updateImage(unFuckedKey(e.target.value));
    }, 1000)()
  }

  return e('div', {className: 'col'}, [
    e('div', {key: 'form', className: 'mb-3'}, [
      e('label', {key: 'label', className: 'form-label'}, 'Image Key'),
      e('input', {key: 'input', type: 'text', className: 'form-control', onChange})
    ]),
    image ? e(ImagePanel, {...image, key: image.id}) : e('div', {key: 'empty'})
  ])
}

export default function Form() {
  const [imagePanels, setImagePanels] = useState([null]);

  const appendImagePanel = () => {
    setImagePanels([...imagePanels, null])
  }

  return e('div', {className: 'container-fluid'},
    e('div', {
      key: 'row',
      className: 'row row-cols-4'
    }, imagePanels.map((_, idx) => e(ImagePanelForm, {key: idx, loaded: appendImagePanel}))),
  )
}
