import { firebase } from 'auth';

class API {
  url: string

  constructor(url: string) {
    this.url = url;
  }

  createLadder(name: string): Promise<any> {
    console.log(`${this.url}/ladder`);

    return firebase.getToken().then(token => {
      return fetch(`${this.url}/ladder`, {
        method: 'POST',
        body: JSON.stringify({
          name,
          token
        })
      }).then(res => res.json())
    });
  }
}

const local = new API('http://localhost:8080')

export {
  local
}