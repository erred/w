import { Request, EventType } from "./iglog_pb.js";
import { FollowatchClient } from "./iglog_grpc_web_pb.js";

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
      let options = { authorization: idToken };
      let svc = new FollowatchClient("https://api.seankhliao.com");
      let req = new Request();
      let ul = ``;
      let call = null;

      let p = window.location.pathname.split("/");
      p.push("unkown");
      switch (p[2]) {
        case "events":
          call = svc.eventLog(req, options, (err, res) => {
            if (err) {
              console.log(err);
            }
            ul =
              `<h5>EventLog</h5>` +
              `<ul>${res
                .getEventsList()
                .map(e => `<li>${eventToHTML(e)}</li>`)
                .join("")}</ul><p>done</p>`;
            document.querySelector(".loader").style.display = "none";
            document.querySelector("body").insertAdjacentHTML("beforeend", ul);
          });
          break;
        case "followers":
          call = svc.followers(req, options, (err, res) => {
            if (err) {
              console.log(err);
            }
            ul =
              `<h5>Followers</h5>` +
              `<ul>${res
                .getUsersList()
                .map(u => `<li>${userToHTML(u)}</li>`)
                .join("")}</ul><p>done</p>`;
            document.querySelector(".loader").style.display = "none";
            document.querySelector("body").insertAdjacentHTML("beforeend", ul);
          });
          break;
        case "following":
          call = svc.following(req, options, (err, res) => {
            if (err) {
              console.log(err);
            }
            ul =
              `<h5>Following</h5>` +
              `<ul>${res
                .getUsersList()
                .map(u => `<li>${userToHTML(u)}</li>`)
                .join("")}</ul><p>done</p>`;
            document.querySelector(".loader").style.display = "none";
            document.querySelector("body").insertAdjacentHTML("beforeend", ul);
          });
          break;
        default:
          ul = `
<ul>
  <li><a href="/iglog/events">events</a> | what changed</li>
  <li><a href="/iglog/followers">followers</a> | who's following</li>
  <li><a href="/iglog/following">following</a> | who interests me</li>
</ul>
        `;
          document.querySelector(".loader").style.display = "none";
          document.querySelector("body").insertAdjacentHTML("beforeend", ul);
      }
    })
    .catch(function(error) {
      console.log(error);
    });
}

function userToHTML(u) {
  return `
<a href="https://instagram.com/${u.getUsername()}">@${u.getUsername()}</a>
<mark>${u.getDisplayname()}</mark>
  `;
}
function eventToHTML(e) {
  let type = "unknown";
  switch (e.getType()) {
    case EventType.FOLLOWERGAINED:
      type = "followed you";
      break;
    case EventType.FOLLOWERLOST:
      type = "unfollowed you";
      break;
    case EventType.FOLLOWINGGAINED:
      type = "you followed";
      break;
    case EventType.FOLLOWINGLOST:
      type = "you unfollowed";
      break;
  }
  return `
${userToHTML(e.getUser())}
<br>
<time datetime="${e.getTime()}">${e.getTime()}</time><span>${type}</span>`;
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
