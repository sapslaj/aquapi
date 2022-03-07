const { createElement: e, useState } = React;

import graphqlQuery from '../graphqlQuery.js';

const addTagMutation = `
mutation($imageId: String!, $tag: String!) {
  AddTags(input: {imageId: $imageId, tags: [$tag]}) {
    id
    url
    tags
  }
}
`;

const removeTagMutation = `
mutation($imageId: String!, $tag: String!) {
  RemoveTags(input: {imageId: $imageId, tags: [$tag]}) {
    id
    url
    tags
  }
}
`;

export default function ImageTagToggle({updateImageTag, id, tag, value}) {
  const [checked, setChecked] = useState(value);
  const [disabled, setDisabled] = useState(false);

  function onChange(event) {
    setChecked(event.target.checked)
    setDisabled(true)
    const variables = {imageId: id, tag};
    const query = event.target.checked ? addTagMutation : removeTagMutation;
    graphqlQuery({query, variables})
      .then(r => r.json())
        .then((result) => {
          if ('errors' in result) {
            event.target.indeterminate = true;
            console.log(result);
          } else {
            updateImageTag(tag, event.target.checked)
          }
          setDisabled(false);
        }, (error) => {
          event.target.indeterminate = true;
          setDisabled(false);
          console.log(error);
        });
    }

  return e('div', {className: 'form-check'}, [
    e('input', {
      key: `${id}-${tag}-checkbox`,
      className: 'form-check-input',
      type: 'checkbox',
      disabled,
      checked,
      onChange,
    }),
    e('div', {key: `${id}-${tag}-checkbox-label`, className: 'form-check-label'}, tag)
  ])
}
