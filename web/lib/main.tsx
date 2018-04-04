import * as React from 'react';
import { ReactiveComponent, Computation, observable } from 'telekinetic';
import { render } from 'react-dom';
import { grpc } from 'grpc-web-client';

import { WikiTribeService } from './_proto/wikitribe_pb_service';
import { UserResponse, CreateUserRequest, ListUsersRequest } from './_proto/wikitribe_pb';

// @ts-ignore
import * as styles from './main.scss';

class UserForm {
  @observable
  name = ''
}

class User {
  id: string
  name: string

  constructor(id: string, name: string) {
    this.id = id;
    this.name = name;
  }
}

const RPC_HOST = 'http://localhost:3000';

class App extends ReactiveComponent {
  @observable
  message = '';

  @observable
  newUser = new UserForm();

  @observable
  users = new Array<User>();

  construct(comp: Computation) {
    const request = new ListUsersRequest();
    grpc.invoke(WikiTribeService.ListUsers, {
      request: request,
      host: RPC_HOST,
      onMessage: (user: UserResponse) => {
        this.users = this.users.concat([
          new User(user.getId(), user.getName())]);
      },
      onEnd: () => { }
    });
  }

  compute() {
    return (
      <div className={styles.test}>
        <h2>Create user</h2>
        <input value={this.newUser.name}
          onChange={(e) => this.newUser.name = e.target.value} />
        <button onClick={() => this.createUser()}>Create</button>
        <h2>Users</h2>
        {this.users.map((user) =>
          <div key={user.id}>{user.name}</div>
        )}
      </div>
    );
  }

  createUser() {
    const req = new CreateUserRequest();
    req.setName(this.newUser.name);
    grpc.invoke(WikiTribeService.CreateUser, {
      request: req,
      host: RPC_HOST,
      onMessage: (user: UserResponse) => {
        console.log(user);
      },
      onEnd: () => { }
    });
  }
}

render(<App />, document.getElementById('root'));
