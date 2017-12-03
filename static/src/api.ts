import { firebase } from 'auth';

class API {
  url: string

  constructor(url: string) {
    this.url = url;
  }

  // attach the firebase token header to the request, return the response as JSON
  private static fetchWithToken(url: string, init?: RequestInit): Promise<any> {
    return firebase.getToken()
      .then(token => fetch(url, Object.assign({}, init, { headers: { "Firebase-Token": token }})))
      .then(res => res.json());
  }

  createLadder(name: string): Promise<any> {
    return API.fetchWithToken(`${this.url}/ladder`, {
      method: 'POST',
      body: JSON.stringify({ 
        name 
      })
    });
  }

  getLadders(): Promise<any> {
    return API.fetchWithToken(`${this.url}/ladders`);
  }

  registerPlayer(): Promise<any> {
    return API.fetchWithToken(`${this.url}/player`, {
      method: 'POST'
    })
  }
}

export default new API(API_URL);