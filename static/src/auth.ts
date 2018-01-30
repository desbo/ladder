import * as firebase from 'firebase';
import { Store } from 'redux';
import { User } from 'firebase';
import API from 'api';

let firebaseApp: null | FirebaseApp = null

export class FirebaseApp {
  app: firebase.app.App;
  
  constructor(config: Object) {
    this.app = firebase.initializeApp(config);
  }

  register(username: string, email: string, password: string): Promise<any> {
    try {
      return this.app.auth().createUserWithEmailAndPassword(email, password)
        .then(() => this.updateProfile(username, ''))
        .then(() => API.registerPlayer())
        .then(() => this.app.auth().currentUser);
    } catch (e) {
      return Promise.reject(e);
    }
  }

  signIn(email: string, password: string): Promise<User> {
    return this.app.auth().signInWithEmailAndPassword(email, password)
      .then(() => this.app.auth().currentUser);;
  }

  updateProfile(displayName: string, photoURL: string): Promise<any> {
    return this.app.auth().currentUser.updateProfile({
      displayName,
      photoURL
    })
  }

  signOut(): Promise<any> {
    return this.app.auth().signOut();
  }

  currentUser(): firebase.User | null {
    return this.app.auth().currentUser;
  }

  getToken(): Promise<any> {
    if (this.currentUser()) 
      return this.currentUser().getIdToken();
    else
      return Promise.resolve();
  }

  signedIn(): boolean {
    return !!this.currentUser();
  }
}

export function initFirebase(): FirebaseApp {
  if (!firebaseApp) firebaseApp = new FirebaseApp(FIREBASE_CONFIG);
  return firebaseApp;
}

export { 
  firebaseApp as firebase
}