import { firebase } from 'auth';

class API {
  url: string

  constructor(url: string) {
    this.url = url;
  }

  private static fetchWithToken(url: string, init?: RequestInit) {
    return firebase.getToken().then(token => 
      fetch(url, Object.assign({}, init, {
        headers: {
          "Firebase-Token": token
        }
      })));
  }

  createLadder(name: string): Promise<any> {
    return API.fetchWithToken(`${this.url}/ladder`, {
      method: 'POST',
      body: JSON.stringify({ 
        name
      })
    }).then(res => res.json())
  }

  registerPlayer(): Promise<any> {
    return API.fetchWithToken(`${this.url}/player`, {
      method: 'POST'
    })
  }
}

const local = new API('http://localhost:8080');

export {
  local
}