import { Msg } from "./authed_pb.js";
import { AuthedClient } from "./authed_grpc_web_pb.js";

import * as firebase from "firebase/app";
import * as firebaseui from "firebaseui";

import "firebase/auth";

const firebaseConfig = {
  apiKey: "AIzaSyAZwB-8GDcap51t7cDUm1BDe3wN3f-DS3o",
  authDomain: "com-seankhliao.firebaseapp.com",
  databaseURL: "https://com-seankhliao.firebaseio.com",
  projectId: "com-seankhliao",
  storageBucket: "com-seankhliao.appspot.com",
  messagingSenderId: "330311169810",
  appId: "1:330311169810:web:6f914fab94f0b716"
};

window.addEventListener("load", () => {
  firebase.initializeApp(firebaseConfig);
  firebase.auth().onAuthStateChanged(user => (user ? signedIn(user) : signedOut()));
});

function signedIn(user) {
  document.querySelector("#firebaseui-auth-container").style.display = "none";
  console.log("onstatechanged signed in");
  document.querySelector(".loader").style.display = "block";

  firebase
    .auth()
    .currentUser.getIdToken(/* forceRefresh */ true)
    .then(function(idToken) {
      let msg = new Msg();
      msg.setMsg("hello, world");

      let options = { authorization: idToken };

      let svc = new AuthedClient("https://api.seankhliao.com");
      let call = svc.echo(msg, options, (err, res) => {
        if (err) {
          console.log(err);
          return;
        }
        document.querySelector(".loader").style.display = "none";
        document.querySelector("body").insertAdjacentHTML("beforeend", `<p>${res.getMsg()}</p>`);
      });
    })
    .catch(function(error) {
      console.log(error);
    });
}

function signedOut() {
  let ui = new firebaseui.auth.AuthUI(firebase.auth());
  let uiConfig = {
    callbacks: {
      signInSuccessWithAuthResult: function(authResult, redirectUrl) {
        console.log(authResult, redirectUrl);
        return true;
      },
      uiShown: function() {
        document.querySelector(".loader").style.display = "none";
      }
    },
    signInFlow: "popup",
    signInSuccessUrl: "/authed",
    signInOptions: [firebase.auth.GoogleAuthProvider.PROVIDER_ID],
    tosUrl: "/terms",
    privacyPolicyUrl: "/privacy"
  };
  ui.start("#firebaseui-auth-container", uiConfig);
}
