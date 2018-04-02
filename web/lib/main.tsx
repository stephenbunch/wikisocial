import * as React from 'react';
import { ReactiveComponent, Computation, observable } from 'telekinetic';
import { render } from 'react-dom';
import { grpc } from 'grpc-web-client';

import { PostService } from './_proto/post_pb_service';
import { Post, GetPostRequest } from './_proto/post_pb';

// @ts-ignore
import * as styles from './main.scss';

class App extends ReactiveComponent {
  @observable
  message = '';

  construct(comp: Computation) {
    const request = new GetPostRequest();
    request.setId('42');
    grpc.invoke(PostService.GetPost, {
      request: request,
      host: 'http://localhost:3000',
      onMessage: (post: Post) => {
        this.message = post.getMessage();
      },
      onEnd: () => { }
    });
  }

  compute(props: {}) {
    return (
      <div className={styles.test}>
        {this.message}
      </div>
    );
  }
}

render(<App />, document.getElementById('root'));
