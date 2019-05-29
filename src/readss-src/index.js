import { ListRequest, ListReply } from "./readss_pb.js";
import { ListerClient } from "./readss_grpc_web_pb.js";

window.addEventListener("load", () => {
  let svc = new ListerClient("http://localhost:8090");

  let call = svc.list(new ListRequest(), null, (err, res) => {
    if (err) {
      console.log(err);
      return;
    }
    let df = document.createDocumentFragment();
    res.getArticlesList().forEach(a => {
      let e = document.createElement("li");
      console.log(a);
      e.innerHTML = `<p><a href="${a.getUrl()}">${a.getTitle()}</a></p><p><time datetime="${a.getTime()}">${a.getTime()}</time> ${a.getReltime()} | <mark>${a.getSource()}</mark></p>`;
      df.appendChild(e);
    });
    document.querySelector("#list").appendChild(df);
  });
});
