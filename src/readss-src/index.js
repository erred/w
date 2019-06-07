import { ListRequest, ListReply } from "./readss_pb.js";
import { ListerClient } from "./readss_grpc_web_pb.js";

window.addEventListener("load", () => {
  let svc = new ListerClient("https://api.seankhliao.com");

  let call = svc.list(new ListRequest(), null, (err, res) => {
    if (err) {
      console.log(err);
      return;
    }
    console.log(res.getArticlesList());
    let df = document.createDocumentFragment();
    res.getArticlesList().forEach(a => {
      let e = document.createElement("li");
      e.innerHTML = `
        <a href="${a.getUrl()}">${a.getTitle()}</a>
        <br>
        <time datetime="${a.getTime()}">${a.getTime()}</time> ${a.getReltime()} | <mark>${a.getSource()}</mark>`;
      df.appendChild(e);
    });
    document.querySelector(".sk-cube-grid").remove();
    document.querySelector("#list").appendChild(df);
  });
  call.on("status", s => {
    console.log(s);
  });
});
