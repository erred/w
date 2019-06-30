import { fireConf, uiConf } from "./config.js";

import { Request, EventType } from "./iglog_pb.js";
import { FollowatchClient } from "./iglog_grpc_web_pb.js";

import * as firebase from "firebase/app";
import * as firebaseui from "firebaseui";

import "firebase/auth";

window.addEventListener("load", () => {
  firebase.initializeApp(fireConf);

  let p = window.location.pathname.split("/");
  p.push("unkmown");
  switch (p[2]) {
    case "events":
    case "followers":
    case "following":
      firebase.auth().onAuthStateChanged(user => (user ? signedIn(user) : signedOut()));
      break;
    default:
      showDefault();
  }
});

function signedIn(user) {
  document.querySelector("#firebaseui-auth-container").style.display = "none";
  document.querySelector(".loader").style.display = "block";

  document.querySelector("footer").insertAdjacentHTML("beforeend", `| <a id="logout" href="#">logout</a>`);
  document.querySelector("#logout").addEventListener("click", e => {
    e.preventDefault();
    firebase
      .auth()
      .signOut()
      .then(() => {
        location.reload();
      })
      .catch(e => {
        console.log("logout", e);
      });
  });

  firebase
    .auth()
    .currentUser.getIdToken()
    .then(idToken => {
      let options = { authorization: idToken };
      let svc = new FollowatchClient("https://api.seankhliao.com");
      let req = new Request();
      let ul = ``;
      let call = null;

      let p = window.location.pathname.split("/");
      switch (p[2]) {
        case "events":
          call = svc.eventLog(req, options, handleShow(showEvents));
          break;
        case "followers":
          call = svc.followers(req, options, handleShow(showUsers));
          break;
        case "following":
          call = svc.following(req, options, handleShow(showUsers));
          break;
      }
    })
    .catch(function(error) {
      console.log(error);
    });
}

function handleShow(handler) {
  return function(err, res) {
    if (err) {
      console.log(err);
      return;
    }
    handler(res);
  };
}
function showEvents(res) {
  showContent(`
    <h5>EventLog</h5>
    <p><span>total: ${res.getEventsList().length}</span></p>
    ${showList(res.getEventsList(), e => `<li>${eventToHTML(e)}</li>`)}`);
}
function showUsers(res) {
  showContent(`
    <h5>Following</h5>
    <p><span>total: ${res.getUsersList().length}</span></p>
    ${showList(res.getUsersList(), u => `<li>${userToHTML(u)}</li>`)}`);
}
function showList(list, lambda) {
  return `
    <ul>
      ${list.map(lambda).join("")}
    </ul>
  `;
}
function showContent(ul) {
  document.querySelector(".loader").style.display = "none";
  document.querySelector("main").insertAdjacentHTML("beforeend", ul);
}

function userToHTML(u) {
  return `
<a href="https://instagram.com/${u.getUsername()}">@${u.getUsername()}</a>
<em>${u.getDisplayname()}</em>
  `;
}

function eventToHTML(e) {
  let type = "unknown";
  switch (e.getType()) {
    case EventType.FOLLOWERGAINED:
      type = '<span class="g">+follower </span>';
      break;
    case EventType.FOLLOWERLOST:
      type = '<span class="r">-follower </span>';
      break;
    case EventType.FOLLOWINGGAINED:
      type = '<span class="g">+following</span>';
      break;
    case EventType.FOLLOWINGLOST:
      type = '<span class="r">-following</span>';
      break;
  }
  return `
${userToHTML(e.getUser())}
<br>
${type} | <time datetime="${e.getTime()}">${e.getTime()}</time>`;
}

function showDefault() {
  showContent(`
    <ul>
      <li><a href="/iglog/events">events</a> | what changed</li>
      <li><a href="/iglog/followers">followers</a> | who's following</li>
      <li><a href="/iglog/following">following</a> | who interests me</li>
    </ul>`);
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
    signInSuccessUrl: window.location.pathname,
    signInOptions: [firebase.auth.GoogleAuthProvider.PROVIDER_ID],
    tosUrl: "/terms",
    privacyPolicyUrl: "/privacy"
  };
  ui.start("#firebaseui-auth-container", uiConfig);
}
