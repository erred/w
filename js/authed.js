import { fireConf, uiConf } from "./config.js";

import { Msg } from "./authed_pb.js";
import { AuthedClient } from "./authed_grpc_web_pb.js";

import * as firebase from "firebase/app";
import * as firebaseui from "firebaseui";

import "firebase/auth";

window.addEventListener("load", () => {
  firebase.initializeApp(fireConf);
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
        let msg = ``;
        if (err) {
          console.log(err);
          msg = `<p>${err}</p>`;
        } else {
          msg = `<p>${res.getMsg()}</p>`;
        }
        document.querySelector(".loader").style.display = "none";
        document.querySelector("main").insertAdjacentHTML("beforeend", msg);
      });
    })
    .catch(function(error) {
      console.log(error);
      document.querySelector(".loader").style.display = "none";
      document.querySelector("main").insertAdjacentHTML("beforeend", `<p>${err}</p>`);
    });
}

function signedOut() {
  let ui = new firebaseui.auth.AuthUI(firebase.auth());
  ui.start("#firebaseui-auth-container", uiConf);
}
