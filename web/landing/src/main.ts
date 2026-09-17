import { courierMarkSource } from "@courier/ui";
import "./app";

const favicon = document.createElement("link");
favicon.rel = "icon";
favicon.type = "image/webp";
favicon.href = courierMarkSource;
document.head.append(favicon);
