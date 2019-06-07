import { ListRequest, ListReply } from "./readss_pb.js";
import { ListerClient } from "./readss_grpc_web_pb.js";

window.addEventListener("load", () => {
  let svc = new ListerClient("https://api.seankhliao.com");

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
    a =>
      `<a href="${a.getUrl()}">${a.getTitle()}</a><br><mark>${a.getSource()}</mark> <time datetime="${a.getTime()}">${a
        .getTime()
        .replace("-", "&#8209;")}&nbsp;${a.getReltime().replace("-", "&#8209;")}</time>`
  )
  .join("")}
</ul>
`;

    document.querySelector(".sk-cube-grid").remove();
    document.querySelector("body").insertAdjacentHTML("beforeend", ul);
  });
  call.on("status", s => {
    console.log(s);
  });
});
