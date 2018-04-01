import * as React from 'react';
import { ReactiveComponent, Computation } from 'telekinetic';
import { render } from 'react-dom';

import * as styles from './main.scss';

class App extends ReactiveComponent {
  construct(comp: Computation) { }
  compute(props: {}) {
    return <div className={styles.test}>hello world!</div>
  }
}

render(<App />, document.getElementById('root'));
