import React from 'react';
import {marked, Renderer} from 'marked';

// marked renders a bare <table>; kuard's stylesheet expects the Bootstrap classes. Delegating to
// the base renderer and adding the attribute keeps us out of the token-to-HTML business, which is
// what changed shape between marked 0.3 and 16.
class TableRenderer extends Renderer {
  table(token) {
    return super.table(token).replace('<table>', '<table class="table table-condensed table-bordered">');
  }
}

marked.setOptions({
  renderer: new TableRenderer(),
  gfm: true,
  breaks: false,
  pedantic: false
});

export default class MarkdownElement extends React.Component {
  render() {
    const { text } = this.props,
    html = marked.parse(text || '');

    return (
      <div>
        <div dangerouslySetInnerHTML={{__html: html}} />
      </div>
    );
  }
}

MarkdownElement.defaultProps = {
  text: ''
};
