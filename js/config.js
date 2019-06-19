import * as firebase from "firebase/app";
import "firebase/auth";

const fireConf = {
  apiKey: "AIzaSyAZwB-8GDcap51t7cDUm1BDe3wN3f-DS3o",
  authDomain: "com-seankhliao.firebaseapp.com",
  databaseURL: "https://com-seankhliao.firebaseio.com",
  projectId: "com-seankhliao",
  storageBucket: "com-seankhliao.appspot.com",
  messagingSenderId: "330311169810",
  appId: "1:330311169810:web:6f914fab94f0b716"
};

var uiConf = {
  callbacks: {
    uiShown: function() {
      document.querySelector(".loader").style.display = "none";
    }
  },
  signInFlow: "popup",
  signInSuccessUrl: window.location.pathname,
  signInOptions: [firebase.auth.GoogleAuthProvider.PROVIDER_ID],
  tosUrl: "/terms",
  privacyPolicyUrl: "/privacy"
};

export { fireConf, uiConf };
