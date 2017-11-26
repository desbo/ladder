import * as firebase from 'firebase';
import { Store } from 'redux';

const config = {
  apiKey: 'AIzaSyBeUnA3-wodE50jIpaEqOVA_h_SXmxCOOQ',
  authDomain: 'tt-ladder.firebaseapp.com',
  databaseURL: 'https://tt-ladder.firebaseio.com',
  projectId: 'tt-ladder',
  storageBucket: 'tt-ladder.appspot.com',
  messagingSenderId: '230896970422'
};

let firebaseApp: null | FirebaseApp = null

export class FirebaseApp {
  app: firebase.app.App;
  
  constructor(config: Object) {
    this.app = firebase.initializeApp(config);
  }

  register(username: string, email: string, password: string): Promise<any> {
    return this.app.auth().createUserWithEmailAndPassword(email, password);
  }

  signIn(email: string, password: string): Promise<any> {
    return this.app.auth().signInWithEmailAndPassword(email, password);
  }

  signOut(): Promise<any> {
    return this.app.auth().signOut();
  }

  currentUser(): firebase.User | null {
    return this.app.auth().currentUser;
  }

  signedIn(): boolean {
    return !!this.currentUser();
  }
}

export function initFirebase(): FirebaseApp {
  if (!firebaseApp) firebaseApp = new FirebaseApp(config);
  return firebaseApp;
}

export { 
  firebaseApp as firebase
}