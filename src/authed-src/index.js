import { Msg } from "./authed_pb.js";
import { authedClient } from "./authed_grpc_web_pb.js";

window.addEventListener("load", () => {
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

      let svc = new authedClient("https://api.seankhliao.com");
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
    tosUrl: "/terms-of-service",
    privacyPolicyUrl: "/privacy"
  };
  ui.start("#firebaseui-auth-container", uiConfig);
}
