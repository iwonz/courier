import { relayMarkSource } from "@courier/ui";
import "./app";

const favicon = document.createElement("link");
favicon.rel = "icon";
favicon.type = "image/png";
favicon.href = relayMarkSource;
document.head.append(favicon);
