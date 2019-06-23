import { ListRequest, ListReply } from "./readss_pb.js";
import { ListerClient } from "./readss_grpc_web_pb.js";

var svc = new ListerClient("https://api.seankhliao.com");
window.addEventListener("load", () => {
  let call = svc.list(new ListRequest(), null, (err, res) => {
    if (err) {
      console.log(err);
      return;
    }
    let ul = `
<ul>
${res
  .getArticlesList()
  .map(
    a => `
    <li>
      <a href="${a.getUrl()}">${a.getTitle()}</a>
      <br>
      <mark>${a.getSource()}</mark>
      <time datetime="${a.getTime()}">${a.getTime()} ${a.getReltime()}</time>
    </li>
    `
  )
  .join("")}
</ul>
`;

    document.querySelector("main").insertAdjacentHTML("beforeend", ul);
    let l = document.querySelector(".loader");
    if (l) {
      l.remove();
    }
  });
  // call.on("status", s => {
  //   console.log(s);
  // });
});
