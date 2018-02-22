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

  getLadder(id: string): Promise<Ladder> {
    return API.fetchWithToken(`${this.url}/ladder/${id}`)
  }

  getLadders(): Promise<LadderState> {
    return API.fetchWithToken(`${this.url}/ladders`);
  }

  joinLadder(id: string): Promise<any> {
    return API.fetchWithToken(`${this.url}/join/${id}`, {
      method: 'POST'
    })
  }

  registerPlayer(): Promise<any> {
    return API.fetchWithToken(`${this.url}/player`, {
      method: 'POST',
      body: JSON.stringify({
        name: firebase.currentUser().displayName 
      })
    })
  }

  submitGame(ladderID: string, opponent: LadderPlayer, myScore: number, theirScore: number) {
    return API.fetchWithToken(`${this.url}/game`, {
      method: 'POST',
      body: JSON.stringify({
        opponentKey: opponent.key,
        ladderID,
        myScore,
        theirScore
      })
    })
  }

  getChart(id: string): Promise<ChartData> {
    return API.fetchWithToken(`${this.url}/chart/${id}`)
  }

  startNewSeason(ladderID: string): Promise<any> {
    return API.fetchWithToken(`${this.url}/ladder/${ladderID}/new-season`, {
      method: 'POST'
    })
  }

}

export default new API(API_URL);